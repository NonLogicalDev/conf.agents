#!/usr/bin/env python3
"""Inspect and create small owner plan homes."""

import argparse
import os
import re
import sys
from pathlib import Path
from uuid import UUID

ACTIVE_STATUSES = ("active", "backlog", "completed")
ALL_STATUSES = (*ACTIVE_STATUSES, "archived")
GROUP_PATTERN = re.compile(r"^([A-Za-z][A-Za-z0-9_]*)-(\d+)-.+$")
EXTERNAL_GROUP_PATTERN = re.compile(r"^ext-[A-Za-z][A-Za-z0-9]*-\d+-.+$")
PLAN_PATTERN = re.compile(r"^plan-(\d+)\s+.+\.md$")
GROUP_TYPE_PATTERN = re.compile(r"^[A-Za-z][A-Za-z0-9_]*$")
SLUG_PATTERN = re.compile(r"^[A-Za-z0-9][A-Za-z0-9_-]*$")
OWNER_THREAD_PREFIX = "codex://threads/"
STATUS_PATTERNS = (
    re.compile(r"^\s*(?:-\s*)?status:\s*([A-Za-z_-]+)\b", re.IGNORECASE),
    re.compile(r"^\s*(?:-\s*)?\*\*status:\*\*\s*([A-Za-z_-]+)\b", re.IGNORECASE),
)
STATUS_NAMES = {
    "active": "active",
    "blocked": "active",
    "in-progress": "active",
    "in_progress": "active",
    "backlog": "backlog",
    "planned": "backlog",
    "queued": "backlog",
    "complete": "completed",
    "completed": "completed",
    "done": "completed",
    "abandoned": "archived",
    "archived": "archived",
}


def existing_directory(path: Path) -> None:
    if path.exists() and not path.is_dir():
        raise ValueError(f"expected a directory: {path}")


def validate_group_type(group_type: str) -> str:
    if not GROUP_TYPE_PATTERN.fullmatch(group_type):
        raise ValueError(
            "group type must start with a letter and contain only letters, numbers, or underscores"
        )
    return group_type


def validate_slug(slug: str) -> str:
    if not SLUG_PATTERN.fullmatch(slug):
        raise ValueError("group slug must contain only letters, numbers, hyphens, or underscores")
    return slug


def validate_description(description: str) -> str:
    description = description.strip()
    if description.endswith(".md"):
        description = description[:-3].rstrip()
    if not description or "/" in description or "\\" in description or "\n" in description:
        raise ValueError("plan description must be a nonempty filename without path separators")
    return description


def positive_digits(value: str) -> int:
    try:
        digits = int(value)
    except ValueError as error:
        raise argparse.ArgumentTypeError("digits must be a positive integer") from error
    if digits <= 0:
        raise argparse.ArgumentTypeError("digits must be a positive integer")
    return digits


def next_group_name(plan_root: Path, group_type: str, slug: str, digits: int | None = None) -> str:
    existing_directory(plan_root)
    validate_group_type(group_type)
    validate_slug(slug)
    highest_number = 0
    number_width = 0
    if plan_root.exists():
        for entry in plan_root.iterdir():
            if not entry.is_dir():
                continue
            if EXTERNAL_GROUP_PATTERN.fullmatch(entry.name):
                continue
            match = GROUP_PATTERN.fullmatch(entry.name)
            if match:
                highest_number = max(highest_number, int(match.group(2)))
                number_width = max(number_width, len(match.group(2)))
    width = digits if digits is not None else number_width or 3
    return f"{group_type}-{highest_number + 1:0{width}d}-{slug}"


def next_plan_path(group_root: Path, description: str, status: str) -> Path:
    existing_directory(group_root)
    description = validate_description(description)
    highest_number = 0
    if group_root.exists():
        for status_name in ALL_STATUSES:
            status_directory = group_root / status_name
            if not status_directory.is_dir():
                continue
            for entry in status_directory.iterdir():
                if not entry.is_file():
                    continue
                match = PLAN_PATTERN.fullmatch(entry.name)
                if match:
                    highest_number = max(highest_number, int(match.group(1)))
    return group_root / status / f"plan-{highest_number + 1:02d} {description}.md"


def write_if_missing(path: Path, contents: str) -> None:
    if path.exists():
        if not path.is_file():
            raise ValueError(f"expected a regular file: {path}")
        return
    path.write_text(contents, encoding="utf-8")


def normalize_owner_thread(value: str) -> str:
    thread_id = value.removeprefix(OWNER_THREAD_PREFIX)
    try:
        identifier = UUID(thread_id)
    except (AttributeError, ValueError) as error:
        raise ValueError(
            f"invalid owner thread: {value!r}; expected a UUID or codex://threads/<uuid>"
        ) from error
    return f"{OWNER_THREAD_PREFIX}{identifier}"


def current_owner_thread(owner_thread: str | None, *, required: bool = False) -> str | None:
    configured = normalize_owner_thread(owner_thread) if owner_thread else None
    current_value = os.environ.get("CODEX_THREAD_ID")
    current = normalize_owner_thread(current_value) if current_value else None

    if configured and current and configured != current:
        raise ValueError(
            f"owner thread mismatch: --owner-thread is {configured} but CODEX_THREAD_ID is {current}"
        )
    if required and not (configured or current):
        raise ValueError("owner thread unavailable; set CODEX_THREAD_ID or pass --owner-thread")
    return configured or current


def read_owner_threads(plan_root: Path) -> set[str] | None:
    """Read the allowlist without rewriting its prose or granting membership."""
    path = plan_root / "OWNERS.md"
    if not path.exists() and not path.is_symlink():
        return None
    if not path.is_file():
        raise ValueError(f"expected a regular file: {path}")
    lines = path.read_text(encoding="utf-8").splitlines()
    sections = [index for index, line in enumerate(lines) if line.strip() == "## Owners"]
    if len(sections) != 1:
        raise ValueError(f"invalid owner register: expected one ## Owners section: {path}")
    rows = []
    for line in lines[sections[0] + 1 :]:
        if re.match(r"^#{1,6}\s", line):
            break
        if line.strip():
            rows.append(line.strip())

    def cells(row: str) -> list[str]:
        if not row.startswith("|") or not row.endswith("|"):
            raise ValueError(f"invalid owner register table row: {path}: {row!r}")
        values = [value.strip() for value in row[1:-1].split("|")]
        if len(values) != 3:
            raise ValueError(f"invalid owner register table row: {path}: {row!r}")
        return values

    if len(rows) < 2 or cells(rows[0]) != ["Thread", "Status", "Authorization"]:
        raise ValueError(f"invalid owner register table header: {path}")
    if not all(re.fullmatch(r":?-{3,}:?", value) for value in cells(rows[1])):
        raise ValueError(f"invalid owner register table separator: {path}")
    seen: set[str] = set()
    active: set[str] = set()
    for row in rows[2:]:
        thread, status, authorization = cells(row)
        if not thread.startswith(OWNER_THREAD_PREFIX) or normalize_owner_thread(thread) != thread:
            raise ValueError(f"invalid owner register thread: {path}: {thread!r}")
        if thread in seen:
            raise ValueError(f"duplicate owner register thread: {path}: {thread}")
        if status not in ("active", "revoked"):
            raise ValueError(f"invalid owner register status: {path}: {status!r}")
        if not authorization:
            raise ValueError(f"missing owner register authorization: {path}: {thread}")
        seen.add(thread)
        if status == "active":
            active.add(thread)
    return active


def check_owner_membership(current_thread: str | None, active: set[str]) -> None:
    if not active:
        raise ValueError("owner register has no active owners: OWNERS.md")
    if current_thread and current_thread not in active:
        raise ValueError(
            f"owner thread mismatch: current thread {current_thread} is not active in OWNERS.md"
        )


def owner_register(owner_thread: str, authorization: str) -> str:
    return (
        "# Owner threads\n\n"
        "## Owners\n\n"
        "| Thread | Status | Authorization |\n"
        "| --- | --- | --- |\n"
        f"| {owner_thread} | active | {authorization} |\n"
    )


def root_readme(plan_root: Path, group_name: str) -> str:
    group_heading = "Current Themes" if group_name.startswith("theme-") else "Current Groups"
    return (
        f"# {plan_root.name}\n\n"
        "[Allowed owner threads](OWNERS.md)\n\n"
        f"## {group_heading}\n\n"
        f"- [{group_name}]({group_name}/EXEC_STATE.md)\n"
    )


def project_memory(group_name: str, plan_path: Path) -> str:
    return (
        "# Project Memory\n\n"
        "Read this file first when starting or resuming owner work. Re-dream it from current evidence instead of trusting it blindly.\n\n"
        "## Project Map\n\n"
        "- [Owner entry point](README.md): project overview and top-level records.\n"
        "- [Allowed owner threads](OWNERS.md): verified owner membership and authorization.\n"
        f"- [Current execution]({group_name}/EXEC_STATE.md): current workstream state.\n"
        f"- [Current plan]({plan_path.relative_to(plan_path.parents[2])}): accepted work and verification.\n\n"
        "## Memory Index\n\n"
        "- Add links to durable findings in `_owner/memory/` when they become useful.\n\n"
        "## Audit\n\n"
        "- Last re-dreamed: not yet audited against current project evidence.\n"
        "- Next audit trigger: owner resume, material project-map change, handoff, completion, or 30 minutes of active work.\n"
    )


def execution_state(group_name: str, description: str, status: str) -> str:
    return (
        f"# {group_name}\n\n"
        "## Goal\n\n"
        f"{description}\n\n"
        "## Current State\n\n"
        f"- Status: {status}\n"
        f"- Next action: Continue {description}.\n"
    )


def plan_contents(description: str, status: str) -> str:
    return (
        f"# {description}\n\n"
        f"Status: {status}\n\n"
        "## Purpose\n\n"
        f"{description}\n\n"
        "## Scope\n\n"
        "- Record the accepted work and important exclusions.\n\n"
        "## Completion Criteria\n\n"
        "- [ ] Record the outcome and the evidence that will prove it.\n"
    )


def create_group(
    plan_root: Path,
    group_type: str,
    slug: str,
    description: str,
    status: str,
    digits: int | None = None,
    owner_thread: str | None = None,
) -> Path:
    existing_directory(plan_root)
    description = validate_description(description)
    group_name = next_group_name(plan_root, group_type, slug, digits)
    expected_owner = current_owner_thread(owner_thread, required=True)
    if expected_owner is None:
        raise ValueError("owner thread unavailable")

    readme_path = plan_root / "README.md"
    if (readme_path.exists() or readme_path.is_symlink()) and not readme_path.is_file():
        raise ValueError(f"expected a regular file: {readme_path}")
    active_owners = read_owner_threads(plan_root)
    if active_owners is not None:
        check_owner_membership(expected_owner, active_owners)
    elif plan_root.exists() and any(plan_root.iterdir()):
        # Legacy declarations need explicit migration; they cannot grant membership.
        raise ValueError(f"missing owner register: {plan_root / 'OWNERS.md'}; migrate existing owner records first")

    plan_root.mkdir(parents=True, exist_ok=True)
    write_if_missing(readme_path, root_readme(plan_root, group_name))
    if active_owners is None:
        write_if_missing(
            plan_root / "OWNERS.md",
            owner_register(expected_owner, "Verified calling thread created this new owner root."),
        )

    group_root = plan_root / group_name
    group_root.mkdir()
    (group_root / status).mkdir()

    write_if_missing(group_root / "EXEC_STATE.md", execution_state(group_name, description, status))
    plan_path = next_plan_path(group_root, description, status)
    write_if_missing(plan_path, plan_contents(description, status))
    write_if_missing(plan_root / "MEMORY.md", project_memory(group_name, plan_path))
    return plan_path


def recorded_status(path: Path) -> str | None:
    for line in path.read_text(encoding="utf-8").splitlines():
        for pattern in STATUS_PATTERNS:
            match = pattern.match(line)
            if match:
                return STATUS_NAMES.get(match.group(1).lower())
    return None


def looks_like_group(path: Path) -> bool:
    if GROUP_PATTERN.fullmatch(path.name) or EXTERNAL_GROUP_PATTERN.fullmatch(path.name):
        return True
    if (path / "EXEC_STATE.md").exists():
        return True
    return any((path / status).exists() for status in ALL_STATUSES)


def check_group(group_root: Path, errors: list[str]) -> int:
    if not (group_root / "EXEC_STATE.md").is_file():
        errors.append(f"missing execution state: {group_root / 'EXEC_STATE.md'}")

    for plan in group_root.glob("plan-*.md"):
        errors.append(f"plan must be inside its status directory: {plan}")

    seen_numbers: dict[int, Path] = {}
    plan_count = 0
    for status in ALL_STATUSES:
        status_directory = group_root / status
        if not status_directory.is_dir():
            continue
        for plan in sorted(status_directory.glob("plan-*.md")):
            match = PLAN_PATTERN.fullmatch(plan.name)
            if not match:
                errors.append(f"invalid plan filename: {plan}")
                continue

            plan_count += 1
            number = int(match.group(1))
            if number in seen_numbers:
                errors.append(
                    f"duplicate plan number {number:02d}: {seen_numbers[number]} and {plan}"
                )
            else:
                seen_numbers[number] = plan

            actual_status = recorded_status(plan)
            if actual_status and actual_status != status:
                errors.append(
                    f"status mismatch: {plan} says {actual_status} but is under {status}/"
                )

    if plan_count == 0:
        errors.append(f"missing numbered plan: {group_root}")

    return plan_count


def doctor(plan_root: Path, owner_thread: str | None = None) -> int:
    errors: list[str] = []
    warnings: list[str] = []
    group_count = 0
    plan_count = 0

    if not plan_root.is_dir():
        print(f"ERROR: missing plan root: {plan_root}")
        return 1

    expected_owner: str | None = None
    try:
        expected_owner = current_owner_thread(owner_thread)
    except ValueError as error:
        errors.append(str(error))

    readme_path = plan_root / "README.md"
    if not readme_path.is_file():
        if (plan_root / "_index_.md").is_file():
            warnings.append(
                f"legacy _index_.md found; add README.md when the owner layout is updated: {plan_root}"
            )
        else:
            errors.append(f"missing root README.md: {readme_path}")
    elif any(
        re.match(r"^\s*owner_thread\s*:", line)
        for line in readme_path.read_text(encoding="utf-8", errors="replace").splitlines()
    ):
        warnings.append(
            f"legacy owner_thread in README.md; migrate the declaration to OWNERS.md and remove it from README.md: {readme_path}"
        )

    try:
        active_owners = read_owner_threads(plan_root)
        if active_owners is None:
            errors.append(f"missing owner register: {plan_root / 'OWNERS.md'}; migrate existing owner records first")
        else:
            check_owner_membership(expected_owner, active_owners)
    except (OSError, ValueError) as error:
        errors.append(str(error))

    memory_path = plan_root / "MEMORY.md"
    if readme_path.is_file() and not memory_path.is_file():
        errors.append(f"missing project MEMORY.md: {memory_path}")

    for entry in sorted(plan_root.iterdir()):
        if not entry.is_dir() or entry.name.startswith(("_", ".")):
            continue
        if not looks_like_group(entry):
            continue
        group_count += 1
        plan_count += check_group(entry, errors)

    for warning in warnings:
        print(f"WARNING: {warning}")
    for error in errors:
        print(f"ERROR: {error}")
    print(
        f"Checked {group_count} groups and {plan_count} plans: {len(errors)} errors, {len(warnings)} warnings."
    )
    return 1 if errors else 0


def argument_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(description="Inspect and create owner plan homes.")
    commands = parser.add_subparsers(dest="command", required=True)

    doctor_command = commands.add_parser("doctor", help="check a plan home without changing it")
    doctor_command.add_argument("plan_root", type=Path)
    doctor_command.add_argument(
        "--owner-thread", help="expected owner thread UUID or codex://threads/<uuid>"
    )

    next_group_command = commands.add_parser(
        "next-group", help="print the next theme or explicitly chosen group name"
    )
    next_group_command.add_argument("plan_root", type=Path)
    next_group_command.add_argument(
        "type_or_slug",
        metavar="type-or-slug",
        help="the group type, or the slug when the type defaults to theme",
    )
    next_group_command.add_argument(
        "slug", nargs="?", help="the group slug when a type is provided"
    )
    next_group_command.add_argument(
        "--digits",
        type=positive_digits,
        help="set the group number width; new themes default to three digits",
    )

    next_plan_command = commands.add_parser("next-plan", help="print the next numbered plan path")
    next_plan_command.add_argument("group_root", type=Path)
    next_plan_command.add_argument("description")
    next_plan_command.add_argument("--status", choices=ALL_STATUSES, default="active")

    create_group_command = commands.add_parser(
        "create-group", help="create a root index and one theme or chosen group"
    )
    create_group_command.add_argument("plan_root", type=Path)
    create_group_command.add_argument(
        "type_or_slug",
        metavar="type-or-slug",
        help="the group type, or the slug when the type defaults to theme",
    )
    create_group_command.add_argument(
        "slug_or_description",
        metavar="slug-or-description",
        help="the group slug, or the description when the type defaults to theme",
    )
    create_group_command.add_argument(
        "description", nargs="?", help="the plan description when a type is provided"
    )
    create_group_command.add_argument("--status", choices=("active", "backlog"), default="active")
    create_group_command.add_argument(
        "--digits",
        type=positive_digits,
        help="set the group number width; new themes default to three digits",
    )
    create_group_command.add_argument(
        "--owner-thread", help="owner thread UUID or codex://threads/<uuid>"
    )

    return parser


def main(argv: list[str] | None = None) -> int:
    parser = argument_parser()
    arguments = parser.parse_args(argv)

    try:
        if arguments.command == "doctor":
            return doctor(arguments.plan_root, arguments.owner_thread)
        if arguments.command == "next-group":
            if arguments.slug is None:
                group_type = "theme"
                slug = arguments.type_or_slug
            else:
                group_type = arguments.type_or_slug
                slug = arguments.slug
            print(next_group_name(arguments.plan_root, group_type, slug, arguments.digits))
            return 0
        if arguments.command == "next-plan":
            print(next_plan_path(arguments.group_root, arguments.description, arguments.status))
            return 0
        if arguments.command == "create-group":
            if arguments.description is None:
                group_type = "theme"
                slug = arguments.type_or_slug
                description = arguments.slug_or_description
            else:
                group_type = arguments.type_or_slug
                slug = arguments.slug_or_description
                description = arguments.description
            print(
                create_group(
                    arguments.plan_root,
                    group_type,
                    slug,
                    description,
                    arguments.status,
                    arguments.digits,
                    arguments.owner_thread,
                )
            )
            return 0
    except (OSError, ValueError) as error:
        print(f"error: {error}", file=sys.stderr)
        return 1

    parser.error(f"unknown command: {arguments.command}")
    return 2


if __name__ == "__main__":
    raise SystemExit(main())

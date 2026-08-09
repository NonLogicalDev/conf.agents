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
OWNER_THREAD_PATTERN = re.compile(r"^\s*owner_thread\s*:\s*(.*?)\s*$")
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


def readme_frontmatter(contents: str, readme_path: Path) -> tuple[list[str], int] | None:
    lines = contents.splitlines(keepends=True)
    if not lines or lines[0].rstrip("\r\n") != "---":
        return None

    for index, line in enumerate(lines[1:], start=1):
        if line.rstrip("\r\n") == "---":
            return lines, index
    raise ValueError(f"invalid README frontmatter: missing closing ---: {readme_path}")


def readme_owner_thread(contents: str, readme_path: Path) -> str | None:
    frontmatter = readme_frontmatter(contents, readme_path)
    if frontmatter is None:
        return None

    lines, end = frontmatter
    values = [
        match.group(1) for line in lines[1:end] if (match := OWNER_THREAD_PATTERN.match(line))
    ]
    if not values:
        return None
    if len(values) > 1:
        raise ValueError(f"duplicate owner_thread in root README.md: {readme_path}")

    value = values[0]
    if len(value) >= 2 and value[0] == value[-1] and value[0] in "\"'":
        value = value[1:-1]
    if not value.startswith(OWNER_THREAD_PREFIX):
        raise ValueError(
            f"invalid owner_thread in root README.md: {readme_path}; expected codex://threads/<uuid>"
        )

    try:
        normalized = normalize_owner_thread(value)
    except ValueError as error:
        raise ValueError(
            f"invalid owner_thread in root README.md: {readme_path}: {value!r}"
        ) from error
    if value != normalized:
        raise ValueError(f"invalid owner_thread in root README.md: {readme_path}; use {normalized}")
    return normalized


def add_owner_thread(contents: str, owner_thread: str, readme_path: Path) -> str:
    owner_line = f'owner_thread: "{owner_thread}"\n'
    frontmatter = readme_frontmatter(contents, readme_path)
    if frontmatter is None:
        return f"---\n{owner_line}---\n\n{contents}"

    lines, end = frontmatter
    lines.insert(end, owner_line)
    return "".join(lines)


def root_readme(plan_root: Path, group_name: str, owner_thread: str) -> str:
    return (
        "---\n"
        f'owner_thread: "{owner_thread}"\n'
        "---\n\n"
        f"# {plan_root.name}\n\n"
        "## Current Groups\n\n"
        f"- [{group_name}]({group_name}/EXEC_STATE.md)\n"
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
    existing_readme: str | None = None
    recorded_owner: str | None = None
    if readme_path.exists():
        if not readme_path.is_file():
            raise ValueError(f"expected a regular file: {readme_path}")
        existing_readme = readme_path.read_text(encoding="utf-8")
        recorded_owner = readme_owner_thread(existing_readme, readme_path)
        if recorded_owner and recorded_owner != expected_owner:
            raise ValueError(
                f"owner thread mismatch: {readme_path} belongs to {recorded_owner}, current thread is {expected_owner}"
            )

    plan_root.mkdir(parents=True, exist_ok=True)
    if existing_readme is None:
        write_if_missing(readme_path, root_readme(plan_root, group_name, expected_owner))
    elif recorded_owner is None:
        readme_path.write_text(
            add_owner_thread(existing_readme, expected_owner, readme_path), encoding="utf-8"
        )

    group_root = plan_root / group_name
    group_root.mkdir()
    (group_root / status).mkdir()

    write_if_missing(group_root / "EXEC_STATE.md", execution_state(group_name, description, status))
    plan_path = next_plan_path(group_root, description, status)
    write_if_missing(plan_path, plan_contents(description, status))
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
    else:
        try:
            recorded_owner = readme_owner_thread(
                readme_path.read_text(encoding="utf-8"), readme_path
            )
            if recorded_owner is None:
                errors.append(f"missing owner_thread in root README.md: {readme_path}")
            elif expected_owner and recorded_owner != expected_owner:
                errors.append(
                    f"owner thread mismatch: {readme_path} belongs to {recorded_owner}, current thread is {expected_owner}"
                )
        except ValueError as error:
            errors.append(str(error))

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
        "next-group", help="print the next name for a chosen group type"
    )
    next_group_command.add_argument("plan_root", type=Path)
    next_group_command.add_argument(
        "type_or_slug",
        metavar="type-or-slug",
        help="the group type, or the slug when the type defaults to task",
    )
    next_group_command.add_argument(
        "slug", nargs="?", help="the group slug when a type is provided"
    )
    next_group_command.add_argument(
        "--digits",
        type=positive_digits,
        help="set the group number width; existing widths are inferred by default",
    )

    next_plan_command = commands.add_parser("next-plan", help="print the next numbered plan path")
    next_plan_command.add_argument("group_root", type=Path)
    next_plan_command.add_argument("description")
    next_plan_command.add_argument("--status", choices=ALL_STATUSES, default="active")

    create_group_command = commands.add_parser(
        "create-group", help="create a root index and one plan group"
    )
    create_group_command.add_argument("plan_root", type=Path)
    create_group_command.add_argument(
        "type_or_slug",
        metavar="type-or-slug",
        help="the group type, or the slug when the type defaults to task",
    )
    create_group_command.add_argument(
        "slug_or_description",
        metavar="slug-or-description",
        help="the group slug, or the description when the type defaults to task",
    )
    create_group_command.add_argument(
        "description", nargs="?", help="the plan description when a type is provided"
    )
    create_group_command.add_argument("--status", choices=("active", "backlog"), default="active")
    create_group_command.add_argument(
        "--digits",
        type=positive_digits,
        help="set the group number width; existing widths are inferred by default",
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
                group_type = "task"
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
                group_type = "task"
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

"""Checks for the public owner plan helper."""

import os
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

SCRIPT = Path(__file__).resolve().parents[1] / "scripts" / "plan.py"
OWNER_THREAD_ID = "11111111-1111-4111-8111-111111111111"
OTHER_THREAD_ID = "22222222-2222-4222-8222-222222222222"
OWNER_THREAD = f"codex://threads/{OWNER_THREAD_ID}"
OTHER_THREAD = f"codex://threads/{OTHER_THREAD_ID}"


class PlanHelperTests(unittest.TestCase):
    def setUp(self) -> None:
        self.temporary = tempfile.TemporaryDirectory()
        self.addCleanup(self.temporary.cleanup)
        self.plan_root = Path(self.temporary.name) / "project-plans"

    def run_helper(
        self, *arguments: str | Path, environment: dict[str, str | None] | None = None
    ) -> subprocess.CompletedProcess[str]:
        process_environment = os.environ.copy()
        process_environment["CODEX_THREAD_ID"] = OWNER_THREAD_ID
        for name, value in (environment or {}).items():
            if value is None:
                process_environment.pop(name, None)
            else:
                process_environment[name] = value
        return subprocess.run(
            [sys.executable, str(SCRIPT), *(str(argument) for argument in arguments)],
            check=False,
            capture_output=True,
            text=True,
            env=process_environment,
        )

    def create_group(
        self,
        kind: str = "task",
        slug: str = "first-task",
        description: str = "First task",
        status: str = "active",
    ) -> Path:
        result = self.run_helper(
            "create-group", self.plan_root, kind, slug, description, "--status", status
        )
        self.assertEqual(result.returncode, 0, result.stderr)
        return Path(result.stdout.strip())

    def test_next_group_shares_task_and_epic_numbers(self) -> None:
        self.plan_root.mkdir()
        (self.plan_root / "task-001-prepare").mkdir()
        (self.plan_root / "epic-004-refresh").mkdir()
        (self.plan_root / "ext-DEMO-99-report").mkdir()

        result = self.run_helper("next-group", self.plan_root, "task", "finish")

        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(result.stdout.strip(), "task-005-finish")

    def test_next_group_defaults_to_task_and_preserves_numeric_width(self) -> None:
        self.plan_root.mkdir()
        (self.plan_root / "research-07-investigate").mkdir()
        (self.plan_root / "ext-DEMO-99-report").mkdir()

        result = self.run_helper("next-group", self.plan_root, "finish")

        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(result.stdout.strip(), "task-08-finish")

    def test_next_group_accepts_explicit_digits_on_an_empty_root(self) -> None:
        self.plan_root.mkdir()

        default_type = self.run_helper("next-group", self.plan_root, "finish", "--digits", "2")
        custom_type = self.run_helper(
            "next-group", self.plan_root, "research", "compare", "--digits", "2"
        )

        self.assertEqual(default_type.returncode, 0, default_type.stderr)
        self.assertEqual(default_type.stdout.strip(), "task-01-finish")
        self.assertEqual(custom_type.returncode, 0, custom_type.stderr)
        self.assertEqual(custom_type.stdout.strip(), "research-01-compare")

    def test_next_group_accepts_custom_types_and_shares_local_numbers(self) -> None:
        self.plan_root.mkdir()
        (self.plan_root / "task-001-prepare").mkdir()
        (self.plan_root / "research-007-investigate").mkdir()
        (self.plan_root / "ext-DEMO-99-report").mkdir()

        result = self.run_helper("next-group", self.plan_root, "feature", "add-search")

        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(result.stdout.strip(), "feature-008-add-search")

    def test_next_group_preserves_existing_numeric_width(self) -> None:
        self.plan_root.mkdir()
        (self.plan_root / "research-07-investigate").mkdir()
        (self.plan_root / "feature-09-search").mkdir()
        (self.plan_root / "ext-DEMO-12345-report").mkdir()

        result = self.run_helper("next-group", self.plan_root, "review", "verify-search")

        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(result.stdout.strip(), "review-10-verify-search")

    def test_next_plan_checks_every_status_directory(self) -> None:
        group_root = self.plan_root / "epic-001-refresh"
        for status in ("active", "backlog", "completed", "archived"):
            (group_root / status).mkdir(parents=True)
        (group_root / "active" / "plan-01 Start.md").write_text(
            "Status: active\n", encoding="utf-8"
        )
        (group_root / "completed" / "plan-04 Finish.md").write_text(
            "Status: completed\n", encoding="utf-8"
        )

        result = self.run_helper("next-plan", group_root, "Prepare release", "--status", "backlog")

        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(
            Path(result.stdout.strip()), group_root / "backlog" / "plan-05 Prepare release.md"
        )

    def test_next_plan_follows_an_archived_abandoned_plan(self) -> None:
        plan_path = self.create_group()
        group_root = plan_path.parents[1]
        archived_path = group_root / "archived" / "plan-01 First task.md"
        archived_path.parent.mkdir()
        plan_path.rename(archived_path)
        archived_path.write_text("Status: abandoned\n", encoding="utf-8")

        result = self.run_helper("next-plan", group_root, "Revised approach")

        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(
            Path(result.stdout.strip()), group_root / "active" / "plan-02 Revised approach.md"
        )
        self.assertEqual(archived_path.read_text(encoding="utf-8"), "Status: abandoned\n")

    def test_create_group_creates_only_current_plan_and_preserves_existing_support(self) -> None:
        self.plan_root.mkdir()
        root_readme = self.plan_root / "README.md"
        root_readme.write_text("Keep the existing root index.\n", encoding="utf-8")
        owner_root = self.plan_root / "_owner"
        owner_root.mkdir()
        owner_readme = owner_root / "README.md"
        owner_readme.write_text("Keep the existing owner notes.\n", encoding="utf-8")

        plan_path = self.create_group("epic", "add-search", "Add search", "backlog")

        self.assertEqual(
            root_readme.read_text(encoding="utf-8"),
            f'---\nowner_thread: "{OWNER_THREAD}"\n---\n\nKeep the existing root index.\n',
        )
        self.assertEqual(
            owner_readme.read_text(encoding="utf-8"), "Keep the existing owner notes.\n"
        )
        self.assertEqual(
            plan_path, self.plan_root / "epic-001-add-search" / "backlog" / "plan-01 Add search.md"
        )
        self.assertTrue(plan_path.is_file())

        for folder in ("docs", "artifacts", "handoff", "memory", "tools", "state"):
            self.assertFalse((owner_root / folder).exists(), folder)
        self.assertTrue((self.plan_root / "epic-001-add-search" / "backlog").is_dir())
        self.assertFalse((self.plan_root / "epic-001-add-search" / "active").exists())
        self.assertFalse((self.plan_root / "epic-001-add-search" / "completed").exists())
        self.assertTrue((self.plan_root / "epic-001-add-search" / "EXEC_STATE.md").is_file())
        self.assertEqual(list(owner_root.glob("wart.*.md")), [])

    def test_create_group_accepts_custom_group_type(self) -> None:
        plan_path = self.create_group("research", "compare-options", "Compare options")

        self.assertEqual(
            plan_path,
            self.plan_root
            / "research-001-compare-options"
            / "active"
            / "plan-01 Compare options.md",
        )
        self.assertTrue(plan_path.is_file())

        result = self.run_helper("doctor", self.plan_root)

        self.assertEqual(result.returncode, 0, result.stdout)
        self.assertIn("Checked 1 groups and 1 plans: 0 errors, 0 warnings.", result.stdout)

    def test_create_group_records_owner_thread_frontmatter(self) -> None:
        self.create_group()

        readme = (self.plan_root / "README.md").read_text(encoding="utf-8")

        self.assertTrue(readme.startswith(f'---\nowner_thread: "{OWNER_THREAD}"\n---\n\n'))

    def test_create_group_accepts_explicit_owner_without_environment(self) -> None:
        result = self.run_helper(
            "create-group",
            self.plan_root,
            "first-task",
            "First task",
            "--owner-thread",
            OTHER_THREAD_ID,
            environment={"CODEX_THREAD_ID": None},
        )

        self.assertEqual(result.returncode, 0, result.stderr)
        readme = (self.plan_root / "README.md").read_text(encoding="utf-8")
        self.assertIn(f'owner_thread: "{OTHER_THREAD}"', readme)

    def test_create_group_accepts_explicit_owner_uri_without_environment(self) -> None:
        result = self.run_helper(
            "create-group",
            self.plan_root,
            "first-task",
            "First task",
            "--owner-thread",
            OTHER_THREAD,
            environment={"CODEX_THREAD_ID": None},
        )

        self.assertEqual(result.returncode, 0, result.stderr)
        readme = (self.plan_root / "README.md").read_text(encoding="utf-8")
        self.assertIn(f'owner_thread: "{OTHER_THREAD}"', readme)

    def test_create_group_rejects_conflicting_explicit_owner(self) -> None:
        result = self.run_helper(
            "create-group",
            self.plan_root,
            "first-task",
            "First task",
            "--owner-thread",
            OTHER_THREAD,
        )

        self.assertEqual(result.returncode, 1)
        self.assertIn("owner thread mismatch", result.stderr)
        self.assertFalse(self.plan_root.exists())

    def test_create_group_requires_a_verified_owner(self) -> None:
        result = self.run_helper(
            "create-group",
            self.plan_root,
            "first-task",
            "First task",
            environment={"CODEX_THREAD_ID": None},
        )

        self.assertEqual(result.returncode, 1)
        self.assertIn("owner thread unavailable", result.stderr)
        self.assertFalse(self.plan_root.exists())

    def test_create_group_preserves_existing_frontmatter_and_body(self) -> None:
        self.plan_root.mkdir()
        readme = self.plan_root / "README.md"
        readme.write_text("---\ntitle: Cedar\n---\n\n# Keep this document\n", encoding="utf-8")

        self.create_group()

        self.assertEqual(
            readme.read_text(encoding="utf-8"),
            f'---\ntitle: Cedar\nowner_thread: "{OWNER_THREAD}"\n---\n\n# Keep this document\n',
        )

    def test_create_group_preserves_matching_owner(self) -> None:
        self.plan_root.mkdir()
        readme = self.plan_root / "README.md"
        original = f'---\ntitle: Cedar\nowner_thread: "{OWNER_THREAD}"\n---\n\n# Existing owner\n'
        readme.write_text(original, encoding="utf-8")

        self.create_group()

        self.assertEqual(readme.read_text(encoding="utf-8"), original)

    def test_create_group_rejects_another_owner_before_changing_files(self) -> None:
        self.plan_root.mkdir()
        readme = self.plan_root / "README.md"
        original = f'---\nowner_thread: "{OTHER_THREAD}"\n---\n\n# Another owner\n'
        readme.write_text(original, encoding="utf-8")

        result = self.run_helper("create-group", self.plan_root, "first-task", "First task")

        self.assertEqual(result.returncode, 1)
        self.assertIn("owner thread mismatch", result.stderr)
        self.assertIn(OTHER_THREAD, result.stderr)
        self.assertIn(OWNER_THREAD, result.stderr)
        self.assertEqual(readme.read_text(encoding="utf-8"), original)
        self.assertEqual(list(self.plan_root.iterdir()), [readme])

    def test_duplicate_owner_threads_are_rejected_without_changing_files(self) -> None:
        self.plan_root.mkdir()
        readme = self.plan_root / "README.md"
        original = f'---\nowner_thread: "{OWNER_THREAD}"\nowner_thread: "{OTHER_THREAD}"\n---\n\n# Existing owner\n'
        readme.write_text(original, encoding="utf-8")

        create = self.run_helper("create-group", self.plan_root, "first-task", "First task")
        inspect = self.run_helper("doctor", self.plan_root)

        self.assertEqual(create.returncode, 1)
        self.assertIn("duplicate owner_thread", create.stderr)
        self.assertEqual(inspect.returncode, 1)
        self.assertIn("ERROR: duplicate owner_thread", inspect.stdout)
        self.assertEqual(readme.read_text(encoding="utf-8"), original)
        self.assertEqual(list(self.plan_root.iterdir()), [readme])

    def test_create_group_defaults_to_task(self) -> None:
        result = self.run_helper(
            "create-group", self.plan_root, "add-search", "Add search", "--status", "backlog"
        )

        self.assertEqual(result.returncode, 0, result.stderr)
        plan_path = Path(result.stdout.strip())
        self.assertEqual(
            plan_path, self.plan_root / "task-001-add-search" / "backlog" / "plan-01 Add search.md"
        )
        self.assertTrue(plan_path.is_file())
        self.assertFalse((self.plan_root / "_owner").exists())

    def test_create_group_accepts_explicit_digits_on_an_empty_root(self) -> None:
        default_type = self.run_helper(
            "create-group", self.plan_root, "add-search", "Add search", "--digits", "2"
        )

        self.assertEqual(default_type.returncode, 0, default_type.stderr)
        self.assertEqual(
            Path(default_type.stdout.strip()),
            self.plan_root / "task-01-add-search" / "active" / "plan-01 Add search.md",
        )

        custom_root = Path(self.temporary.name) / "custom-plans"
        custom_type = self.run_helper(
            "create-group", custom_root, "research", "compare", "Compare options", "--digits", "2"
        )

        self.assertEqual(custom_type.returncode, 0, custom_type.stderr)
        self.assertEqual(
            Path(custom_type.stdout.strip()),
            custom_root / "research-01-compare" / "active" / "plan-01 Compare options.md",
        )

    def test_create_group_preserves_owner_instructions_legacy_prompt_and_prose(self) -> None:
        self.plan_root.mkdir()
        owner_instructions = self.plan_root / "AGENTS.md"
        legacy_owner_prompt = self.plan_root / "OWNER_PROMPT.md"
        prose_steering = self.plan_root / "PROSE_STEERING.md"
        owner_instructions.write_text("Keep the current owner instructions.\n", encoding="utf-8")
        legacy_owner_prompt.write_text("Keep the older owner guidance.\n", encoding="utf-8")
        prose_steering.write_text("Keep the existing prose guidance.\n", encoding="utf-8")

        plan_path = self.create_group("research", "compare-options", "Compare options")

        self.assertTrue(plan_path.is_file())
        self.assertEqual(
            owner_instructions.read_text(encoding="utf-8"), "Keep the current owner instructions.\n"
        )
        self.assertEqual(
            legacy_owner_prompt.read_text(encoding="utf-8"), "Keep the older owner guidance.\n"
        )
        self.assertEqual(
            prose_steering.read_text(encoding="utf-8"), "Keep the existing prose guidance.\n"
        )
        self.assertIn(
            f'owner_thread: "{OWNER_THREAD}"',
            (self.plan_root / "README.md").read_text(encoding="utf-8"),
        )

    def test_create_group_preserves_another_owners_instructions(self) -> None:
        self.plan_root.mkdir()
        readme = self.plan_root / "README.md"
        readme_contents = f'---\nowner_thread: "{OTHER_THREAD}"\n---\n\n# Another owner\n'
        readme.write_text(readme_contents, encoding="utf-8")
        owner_instructions = self.plan_root / "AGENTS.md"
        owner_instructions.write_text("Keep another owner's instructions.\n", encoding="utf-8")

        result = self.run_helper("create-group", self.plan_root, "first-task", "First task")

        self.assertEqual(result.returncode, 1)
        self.assertIn("owner thread mismatch", result.stderr)
        self.assertEqual(readme.read_text(encoding="utf-8"), readme_contents)
        self.assertEqual(
            owner_instructions.read_text(encoding="utf-8"),
            "Keep another owner's instructions.\n",
        )
        self.assertEqual(set(self.plan_root.iterdir()), {readme, owner_instructions})

    def test_create_group_does_not_create_optional_owner_or_prose_files(self) -> None:
        plan_path = self.create_group()

        self.assertTrue(plan_path.is_file())
        self.assertFalse((self.plan_root / "_owner").exists())
        self.assertFalse((self.plan_root / "AGENTS.md").exists())
        self.assertFalse((self.plan_root / "OWNER_PROMPT.md").exists())
        self.assertFalse((self.plan_root / "PROSE_STEERING.md").exists())
        self.assertFalse((self.plan_root / "CHANGELOG.md").exists())
        self.assertFalse((plan_path.parents[1] / "backlog").exists())
        self.assertFalse((plan_path.parents[1] / "completed").exists())

    def test_create_group_omits_empty_optional_plan_sections(self) -> None:
        plan_path = self.create_group()
        execution_state = plan_path.parents[1] / "EXEC_STATE.md"

        self.assertNotIn("## Workstreams", execution_state.read_text(encoding="utf-8"))
        self.assertNotIn(
            "## Produced External Artifacts", execution_state.read_text(encoding="utf-8")
        )
        self.assertNotIn("## Work Log", plan_path.read_text(encoding="utf-8"))
        self.assertNotIn("## Validation", plan_path.read_text(encoding="utf-8"))

    def test_generated_documents_keep_distinct_responsibilities(self) -> None:
        plan_path = self.create_group()
        root_readme = (self.plan_root / "README.md").read_text(encoding="utf-8")
        execution_state = (plan_path.parents[1] / "EXEC_STATE.md").read_text(encoding="utf-8")
        plan = plan_path.read_text(encoding="utf-8")

        self.assertIn("[task-001-first-task](task-001-first-task/EXEC_STATE.md)", root_readme)
        self.assertNotIn("Next action", root_readme)
        self.assertNotIn("## Purpose", root_readme)

        self.assertIn("- Status: active", execution_state)
        self.assertIn("- Next action: Continue First task.", execution_state)

        self.assertIn("Status: active", plan)
        self.assertIn("## Purpose", plan)
        self.assertIn("## Scope", plan)
        self.assertIn("## Completion Criteria", plan)
        self.assertNotIn("Next action", plan)
        self.assertNotIn("## Current State", plan)

    def test_doctor_accepts_a_new_group(self) -> None:
        self.create_group()

        result = self.run_helper("doctor", self.plan_root)

        self.assertEqual(result.returncode, 0, result.stdout)
        self.assertIn("Checked 1 groups and 1 plans: 0 errors, 0 warnings.", result.stdout)

    def test_doctor_reports_missing_owner_thread(self) -> None:
        self.create_group()
        readme = self.plan_root / "README.md"
        readme.write_text("# Existing owner without metadata\n", encoding="utf-8")

        result = self.run_helper("doctor", self.plan_root)

        self.assertEqual(result.returncode, 1)
        self.assertIn("ERROR: missing owner_thread in root README.md", result.stdout)
        self.assertEqual(readme.read_text(encoding="utf-8"), "# Existing owner without metadata\n")

    def test_doctor_reports_invalid_owner_thread(self) -> None:
        self.create_group()
        readme = self.plan_root / "README.md"
        readme.write_text(
            '---\nowner_thread: "codex://threads/not-a-uuid"\n---\n\n# Existing owner\n',
            encoding="utf-8",
        )

        result = self.run_helper("doctor", self.plan_root)

        self.assertEqual(result.returncode, 1)
        self.assertIn("ERROR: invalid owner_thread in root README.md", result.stdout)

    def test_doctor_reports_owner_thread_mismatch(self) -> None:
        self.create_group()

        result = self.run_helper(
            "doctor", self.plan_root, environment={"CODEX_THREAD_ID": OTHER_THREAD_ID}
        )

        self.assertEqual(result.returncode, 1)
        self.assertIn("ERROR: owner thread mismatch", result.stdout)
        self.assertIn(OWNER_THREAD, result.stdout)
        self.assertIn(OTHER_THREAD, result.stdout)

    def test_doctor_accepts_explicit_owner_without_environment(self) -> None:
        self.create_group()

        result = self.run_helper(
            "doctor",
            self.plan_root,
            "--owner-thread",
            OWNER_THREAD,
            environment={"CODEX_THREAD_ID": None},
        )

        self.assertEqual(result.returncode, 0, result.stdout)

    def test_doctor_checks_metadata_without_a_current_thread(self) -> None:
        self.create_group()

        result = self.run_helper("doctor", self.plan_root, environment={"CODEX_THREAD_ID": None})

        self.assertEqual(result.returncode, 0, result.stdout)

    def test_doctor_rejects_explicit_owner_conflicting_with_current_thread(self) -> None:
        self.create_group()

        result = self.run_helper("doctor", self.plan_root, "--owner-thread", OTHER_THREAD)

        self.assertEqual(result.returncode, 1)
        self.assertIn("ERROR: owner thread mismatch", result.stdout)

    def test_doctor_reports_missing_readme_execution_state_and_plan(self) -> None:
        self.plan_root.mkdir()
        (self.plan_root / "task-001-broken").mkdir()

        result = self.run_helper("doctor", self.plan_root)

        self.assertEqual(result.returncode, 1)
        self.assertIn("missing root README.md", result.stdout)
        self.assertIn("missing execution state", result.stdout)
        self.assertIn("missing numbered plan", result.stdout)
        self.assertNotIn("missing status directory", result.stdout)

    def test_doctor_accepts_partial_optional_owner_support(self) -> None:
        self.create_group()
        (self.plan_root / "_owner" / "memory").mkdir(parents=True)

        result = self.run_helper("doctor", self.plan_root)

        self.assertEqual(result.returncode, 0, result.stdout)
        self.assertIn("Checked 1 groups and 1 plans: 0 errors, 0 warnings.", result.stdout)

    def test_doctor_ignores_unrelated_optional_root_directories(self) -> None:
        self.create_group()
        (self.plan_root / "scratch").mkdir()

        result = self.run_helper("doctor", self.plan_root)

        self.assertEqual(result.returncode, 0, result.stdout)
        self.assertIn("Checked 1 groups and 1 plans: 0 errors, 0 warnings.", result.stdout)

    def test_doctor_treats_a_legacy_index_as_a_warning(self) -> None:
        self.plan_root.mkdir()
        (self.plan_root / "_index_.md").write_text("# Existing index\n", encoding="utf-8")

        result = self.run_helper("doctor", self.plan_root)

        self.assertEqual(result.returncode, 0, result.stdout)
        self.assertIn("WARNING: legacy _index_.md found", result.stdout)
        self.assertNotIn("ERROR: missing root README.md", result.stdout)

    def test_doctor_reports_duplicate_plan_numbers(self) -> None:
        plan_path = self.create_group()
        duplicate = plan_path.parents[1] / "completed" / "plan-01 Another task.md"
        duplicate.parent.mkdir()
        duplicate.write_text("Status: completed\n", encoding="utf-8")

        result = self.run_helper("doctor", self.plan_root)

        self.assertEqual(result.returncode, 1)
        self.assertIn("duplicate plan number 01", result.stdout)

    def test_doctor_reports_a_status_mismatch(self) -> None:
        plan_path = self.create_group()
        plan_path.write_text("# First task\n\n- Status: completed\n", encoding="utf-8")

        result = self.run_helper("doctor", self.plan_root)

        self.assertEqual(result.returncode, 1)
        self.assertIn("status mismatch", result.stdout)
        self.assertIn("says completed but is under active/", result.stdout)

    def test_doctor_accepts_an_abandoned_plan_in_archived(self) -> None:
        plan_path = self.create_group()
        archived_path = plan_path.parents[1] / "archived" / plan_path.name
        archived_path.parent.mkdir()
        plan_path.rename(archived_path)
        archived_path.write_text("# First task\n\n- Status: abandoned\n", encoding="utf-8")

        result = self.run_helper("doctor", self.plan_root)

        self.assertEqual(result.returncode, 0, result.stdout)
        self.assertIn("Checked 1 groups and 1 plans: 0 errors, 0 warnings.", result.stdout)

    def test_doctor_rejects_an_abandoned_plan_in_active(self) -> None:
        plan_path = self.create_group()
        plan_path.write_text("# First task\n\n- Status: abandoned\n", encoding="utf-8")

        result = self.run_helper("doctor", self.plan_root)

        self.assertEqual(result.returncode, 1)
        self.assertIn("status mismatch", result.stdout)
        self.assertIn("under active/", result.stdout)

    def test_doctor_allows_a_blocked_plan_in_active(self) -> None:
        plan_path = self.create_group()
        plan_path.write_text("# First task\n\n- Status: blocked\n", encoding="utf-8")

        result = self.run_helper("doctor", self.plan_root)

        self.assertEqual(result.returncode, 0, result.stdout)

    def test_invalid_slug_does_not_create_files(self) -> None:
        result = self.run_helper("create-group", self.plan_root, "task", "bad/slug", "First task")

        self.assertEqual(result.returncode, 1)
        self.assertIn("group slug must contain only", result.stderr)
        self.assertFalse(self.plan_root.exists())

    def test_invalid_group_type_does_not_create_files(self) -> None:
        for group_type in ("bad/type", "bad-type", "../escape", "1feature"):
            with self.subTest(group_type=group_type):
                result = self.run_helper(
                    "create-group", self.plan_root, group_type, "first-task", "First task"
                )

                self.assertEqual(result.returncode, 1)
                self.assertIn("group type must start with a letter", result.stderr)
                self.assertFalse(self.plan_root.exists())

    def test_nonpositive_digits_do_not_create_files(self) -> None:
        for digits in ("0", "-1"):
            with self.subTest(command="next-group", digits=digits):
                result = self.run_helper(
                    "next-group", self.plan_root, "first-task", "--digits", digits
                )
                self.assertEqual(result.returncode, 2)
                self.assertIn("digits must be a positive integer", result.stderr)
                self.assertFalse(self.plan_root.exists())

            with self.subTest(command="create-group", digits=digits):
                result = self.run_helper(
                    "create-group", self.plan_root, "first-task", "First task", "--digits", digits
                )
                self.assertEqual(result.returncode, 2)
                self.assertIn("digits must be a positive integer", result.stderr)
                self.assertFalse(self.plan_root.exists())


if __name__ == "__main__":
    unittest.main()

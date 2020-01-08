from pathlib import Path
import re
from zipfile import ZipFile

from invoke import task

# Configuration values.
BUILD_DIR = "dist"
GOARCH = "amd64"
PLATFORMS = ["linux", "darwin", "windows"]
PROJECT = "keeparelease"
TEST_CMD = f"go test -v -cover -coverprofile=coverage.out  ./{PROJECT}"


@task
def ci(c):
    """Run the full CI suite"""
    c.run("goimports -d .")
    c.run("golint ./...")
    c.run("go vet ./...")
    c.run(TEST_CMD)


@task
def clean(c):
    """Remove unwanted files in project (!DESTRUCTIVE!)."""
    c.run("git clean -ffdx")
    c.run("git reset --hard")


@task
def dist(c, p=""):
    """Build binaries for the targeted platforms."""
    build_dir = Path(BUILD_DIR).resolve()

    # Get the tag.
    tag = c.run("git describe", hide=True)

    # Build all the projects.
    p = Path(PROJECT)
    output = build_dir / p.name
    builder(c, p, tag.stdout.strip(), output)


@task
def publish(c):
    """Create a GitHub release."""
    build_dir = Path(BUILD_DIR)
    if not build_dir.exists:
        raise FileNotFoundError("there is nothing to publish.")

    # Get the tag.
    tag = c.run("git describe", hide=True)

    # Get the assets to publish.
    assets = [f"-a {x}" for x in build_dir.iterdir() if x.is_file() and str(x)]

    # Prepare the command to run
    cmd = f'keeparelease -t {tag.stdout.strip()} {" ".join(assets)}'

    # Run it.
    c.run(cmd)


@task(default=True)
def setup(c):
    """Setup the full environment."""
    c.run("go mod tidy")


@task
def test(c):
    """Run the unit tests."""
    c.run(TEST_CMD)


@task
def view_coverage(c, html=False):
    """View code coverage."""
    out = "html" if html else "func"
    c.run(f"go tool cover -{out}=coverage.out")


def builder(c, project, tag, output):
    """Build a project."""
    if not project.exists():
        raise ValueError(f"project {project} cannot be found in {project.resolve()}")
    build_flags = (
        f"-ldflags=\"-X 'github.com/rgreinho/keeparelease/cmd.Version={tag}'\""
    )
    for platform in PLATFORMS:
        cmd = (
            f"GOOS={platform} GOARCH={GOARCH}"
            f" go build {build_flags} -o {output.resolve()}-{tag}-{platform}-{GOARCH}"
        )
        c.run(cmd)


@task(dist, publish)
def release(c):
    """Build and publish."""

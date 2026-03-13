# Release Process

This document describes the release process for `kswp`. We use [semantic versioning](https://semver.org/) for our releases.

## Prerequisites

- You must have a GitHub account with write access to the repository.
- You must have `git` installed and configured on your machine.

## Steps

1.  **Create a release branch:**

    Create a release branch from `master`. The branch name should be `release/vX.Y.Z`, where `X.Y.Z` is the version number of the release.

    ```bash
    git checkout master
    git pull origin master
    git checkout -b release/vX.Y.Z
    ```

2.  **Update the version:**

    Update the version in the following files:

    - `homebrew/kswp.rb`
    - `krew/kswp.yaml`


4.  **Commit the changes:**

    Commit the changes to the release branch.

    ```bash
    git add .
    git commit -m "Prepare for release vX.Y.Z"
    ```

5.  **Push the release branch:**

    Push the release branch to the remote repository.

    ```bash
    git push origin release/vX.Y.Z
    ```

6.  **Create a pull request:**

    Create a pull request from the release branch to the `master` branch.

7.  **Merge the pull request:**

    Once the pull request is approved, merge it into the `master` branch.

8.  **Create a release on GitHub:**

    Create a new release on GitHub from the `master` branch with the tag `vX.Y.Z`. The release notes should contain the same information as the changelog.

9.  **Update the version on the `master` branch:**

    After the release is created, update the version in the following files on the `master` branch to the next development version (e.g., `vX.Y.Z-next`):

    - `homebrew/kswp.rb`
    - `krew/kswp.yaml`

    Commit the changes directly to the `master` branch.

    ```bash
    git checkout master
    git pull origin master
    # Update the version in the files
    git add .
    git commit -m "Prepare for next development version"
    git push origin master
    ```

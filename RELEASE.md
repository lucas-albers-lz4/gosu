# Release Process

Follow these steps to create a new release:

## 1. Increment the Version

Edit `version.go` and update the `Version` constant to the new release version (e.g., `1.22`).

## 2. Commit the Change

```sh
git add version.go
git commit -m "Bump version to v1.22"
```

## 3. Tag the Release

Tag the commit with the new version:

```sh
git tag 1.22
```

## 4. Push the Commit and Tag

```sh
git push
# Push tags
git push origin 1.22
```

## 5. GitHub Actions

A GitHub Actions workflow will automatically build, sign, and publish the release assets to the [GitHub Releases page](https://github.com/lucas-albers-lz4/gosu/releases) for the new tag.

---

**Note:**
- TAG is just a number no leading v so this is valid -> 1.18 and this is not v1.18. Thats we people expect to download the binary so we follow that format
- Binaries are GPG signed automatically using the configured GPG key.
- You can verify signatures using the corresponding `.asc` files and the public key published on [keys.openpgp.org](https://keys.openpgp.org/).

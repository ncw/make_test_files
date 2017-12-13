# Making a release #

Compile and test

Then run

  goreleaser --rm-dist --snapshot

To test the build

When happy, tag the release

  git tag -a v1.0.XX -m "Release v1.0.XX"
  git push --tags origin master

Then do a release build (set GITHUB token first)

  goreleaser --rm-dist

package releasesrepo

import (
	"fmt"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	semiver "github.com/cppforlife/go-semi-semantic/version"

	bhnotesrepo "github.com/cppforlife/bosh-hub/release/notesrepo"
)

type ReleaseVersionRec struct {
	notesRepo bhnotesrepo.NotesRepository

	Source     string
	VersionRaw string
}

type ReleaseVersionRecSorting []ReleaseVersionRec

func (r ReleaseVersionRec) String() string {
	return fmt.Sprintf("%s %s", r.Source, r.VersionRaw)
}

func (r ReleaseVersionRec) SourceShortName() string {
	parts := strings.Split(r.Source, "/")

	return parts[len(parts)-1]
}

// Version returns parsed version
// todo non-memoized lazy loading is expensive
func (r ReleaseVersionRec) Version() semiver.Version {
	// Validate should not allow invalid version to be saved
	ver, err := semiver.NewVersionFromString(r.VersionRaw)
	if err != nil {
		panic(fmt.Sprintf("Version '%s' is not valid: %s", r.VersionRaw, err))
	}

	return ver
}

func (r ReleaseVersionRec) Notes() (bhnotesrepo.NoteRec, bool, error) {
	return r.notesRepo.Find(r.Source, r.VersionRaw)
}

func (r ReleaseVersionRec) SetNotes(noteRec bhnotesrepo.NoteRec) error {
	return r.notesRepo.Save(r.Source, r.VersionRaw, noteRec)
}

func (r ReleaseVersionRec) Validate() error {
	if len(r.Source) == 0 {
		return bosherr.New("Expected source to be non-empty")
	}

	if len(r.VersionRaw) == 0 {
		return bosherr.New("Expected version to be non-empty")
	}

	_, err := semiver.NewVersionFromString(r.VersionRaw)
	if err != nil {
		return bosherr.WrapError(err, "Expected version to be a valid version")
	}

	return nil
}

func (s ReleaseVersionRecSorting) Len() int           { return len(s) }
func (s ReleaseVersionRecSorting) Less(i, j int) bool { return s[i].Version().IsLt(s[j].Version()) }
func (s ReleaseVersionRecSorting) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
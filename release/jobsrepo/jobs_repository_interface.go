package jobsrepo

import (
	bpreljob "github.com/cppforlife/bosh-provisioner/release/job"

	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
)

type JobsRepository interface {
	FindAll(bhrelsrepo.ReleaseVersionRec) ([]bpreljob.Job, error)
	SaveAll(bhrelsrepo.ReleaseVersionRec, []bpreljob.Job) error
}

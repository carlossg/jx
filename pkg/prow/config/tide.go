package config

import (
	"fmt"
	"github.com/jenkins-x/jx/pkg/log"
	"github.com/jenkins-x/jx/pkg/util"
	"k8s.io/test-infra/prow/config"
	"time"
)

// CreateTide creates a default Tide Config object
func CreateTide() config.Tide {
	// todo get the real URL, though we need to handle the multi cluster use case where dev namespace may be another cluster, so pass it in as an arg?
	t := config.Tide{
		TargetURL: "https://tide.foo.bar",
	}

	var qs []config.TideQuery
	qs = append(qs, createApplicationTideQuery())
	qs = append(qs, createEnvironmentTideQuery())
	t.Queries = qs

	myTrue := true
	myFalse := false

	t.SyncPeriod = time.Duration(30)
	t.StatusUpdatePeriod = time.Duration(30)
	t.ContextOptions = config.TideContextPolicyOptions{
		TideContextPolicy: config.TideContextPolicy{
			FromBranchProtection: &myTrue,
			SkipUnknownContexts:  &myFalse,
		},
	}

	return t
}

// AddRepoToTideConfig adds a code repository to the Tide section of the Prow Config
func AddRepoToTideConfig(t *config.Tide, repo string, kind Kind) error {
	switch kind {
	case Application:
		found := false
		for index, q := range t.Queries {
			if util.Contains(q.Labels, "approved") {
				found = true
				repos := t.Queries[index].Repos
				if !util.Contains(repos, repo) {
					repos = append(repos, repo)
					t.Queries[index].Repos = repos
				}
			}
		}

		if !found {
			log.Infof("Failed to find 'application' tide config, adding...\n")
			t.Queries = append(t.Queries, createApplicationTideQuery())
		}
	case Environment:
		found := false
		for index, q := range t.Queries {
			if !util.Contains(q.Labels, "approved") {
				found = true
				repos := t.Queries[index].Repos
				if !util.Contains(repos, repo) {
					repos = append(repos, repo)
					t.Queries[index].Repos = repos
				}
			}
		}

		if !found {
			log.Infof("Failed to find 'environment' tide config, adding...\n")
			t.Queries = append(t.Queries, createEnvironmentTideQuery())
		}
	case Protection:
		// No Tide config needed for Protection
	default:
		return fmt.Errorf("unknown Prow config kind %s", kind)
	}
	return nil
}

// RemoveRepoFromTideConfig adds a code repository to the Tide section of the Prow Config
func RemoveRepoFromTideConfig(t *config.Tide, repo string, kind Kind) error {
	switch kind {
	case Application:
		found := false
		for index, q := range t.Queries {
			if util.Contains(q.Labels, "approved") {
				found = true
				t.Queries[index].Repos = util.RemoveStringFromSlice(t.Queries[index].Repos, repo)
			}
		}

		if !found {
			log.Infof("Failed to find 'application' tide config, adding...\n")
		}
	case Environment:
		found := false
		for index, q := range t.Queries {
			if !util.Contains(q.Labels, "approved") {
				found = true
				t.Queries[index].Repos = util.RemoveStringFromSlice(t.Queries[index].Repos, repo)
			}
		}

		if !found {
			log.Infof("Failed to find 'environment' tide config, adding...\n")
		}
	case Protection:
		// No Tide config needed for Protection
	default:
		return fmt.Errorf("unknown Prow config kind %s", kind)
	}
	return nil
}

func createApplicationTideQuery() config.TideQuery {
	return config.TideQuery{
		Repos:         []string{"jenkins-x/dummy"},
		Labels:        []string{"approved"},
		MissingLabels: []string{"do-not-merge", "do-not-merge/hold", "do-not-merge/work-in-progress", "needs-ok-to-test", "needs-rebase"},
	}
}

func createEnvironmentTideQuery() config.TideQuery {
	return config.TideQuery{
		Repos:         []string{"jenkins-x/dummy-environment"},
		Labels:        []string{},
		MissingLabels: []string{"do-not-merge", "do-not-merge/hold", "do-not-merge/work-in-progress", "needs-ok-to-test", "needs-rebase"},
	}
}

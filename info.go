package redisimpl

import (
	"context"
	"strings"
)

// based on
// https://redis.io/commands/info/

type Info interface {
	Info(ctx context.Context, sections ...Section) (map[string]interface{}, error)
}

type (
	Section string
)

func (s Section) String() string {
	return string(s)
}

var (
	SectionServer       Section = "server"
	SectionClients      Section = "clients"
	SectionMemory       Section = "memory"
	SectionPersistence  Section = "persistence"
	SectionStats        Section = "stats"
	SectionReplication  Section = "replication"
	SectionCpu          Section = "cpu"
	SectionCommandStats Section = "commandstats"
	SectionLatencyStats Section = "latencystats"
	SectionCluster      Section = "cluster"
	SectionModules      Section = "modules"
	SectionKeyspace     Section = "keyspace"
	SectionErrorStats   Section = "errorstats"

	SectionAll        Section = "all"        // return all sections (excluding module generated ones)
	SectionDefault    Section = "default"    // return only the default set of sections
	SectionEverything Section = "everything" // includes all and modules

	// when no parameter is provided, the default option is assumed
)

func (r *redis) Info(ctx context.Context, sections ...Section) (map[string]interface{}, error) {
	// parse sections
	var sectionStrList []string
	for _, v := range sections {
		sectionStrList = append(sectionStrList, v.String())
	}
	// run cmd info
	res, err := r.client.Info(ctx, sectionStrList...).Result()
	if err != nil {
		return nil, err
	}
	// parse result
	return parseToMap(res), nil
}

func parseToMap(v string) map[string]interface{} {
	maps := make(map[string]interface{}, 0)
	// split line
	for _, line := range strings.Split(v, "\n") {
		// skip comment - has prefix #
		if strings.HasPrefix(line, "#") {
			continue
		}

		// split key-value by ":"
		s := strings.SplitN(line, ":", 2)
		if len(s) != 2 {
			continue
		}

		maps[s[0]] = s[1]
	}
	return maps
}

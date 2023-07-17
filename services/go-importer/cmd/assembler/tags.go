package main

import (
	"go-importer/internal/pkg/db"
	"log"
	"regexp"
)

var flagRegex *regexp.Regexp

func EnsureRegex(reg *string) {
	if flagRegex == nil {
		reg, err := regexp.Compile(*reg)
		if err != nil {
			log.Fatal("Failed to compile flag regex: ", err)
		} else {
			flagRegex = reg
		}
	}
}

func containsTag(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Apply flag in/flag out tags to the entire flow.
// This assumes the `Data` part of the flowItem is already pre-processed, s.t.
// we can run regex tags over the payload directly
func ApplyFlagTags(flow *db.FlowEntry, reg *string) {
	EnsureRegex(reg)

	// If the regex is not valid, bail here
	if flagRegex == nil {
		return
	}

	for idx := 0; idx < len(flow.Flow); idx++ {
		flowItem := &flow.Flow[idx]
		if flagRegex.MatchString(flowItem.Data) {
			var tag string
			if flowItem.From == "c" {
				tag = "flag-in"
			} else {
				tag = "flag-out"
			}
			// Add the tag if it doesn't already exist
			if !containsTag(flow.Tags, tag) {
				flow.Tags = append(flow.Tags, tag)
			}
		}
	}
}

// Apply Sla Tag for packets with source ip of SLA
func ApplySlaTag(flow *db.FlowEntry, sla_ip *string) {

	// If the sla ip is non existent, bail
	if sla_ip == nil {
		return
	}

	if flow.Src_ip == *sla_ip {
		flow.Tags = append(flow.Tags, "sla")
	}
}

// TODO: other tags
// Plan: 
// 		seperate HTTP and regular tcp
// 		add a "libc leak" tag
// 		some others idk
/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package service

import (
	"context"
	"log"
	"time"
)

type AEnvCleaner interface {
	cleanup()
}

type AEnvCleanManager struct {
	cleaner AEnvCleaner

	interval time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewAEnvCleanManager(cleaner AEnvCleaner, duration time.Duration) *AEnvCleanManager {
	ctx, cancel := context.WithCancel(context.Background())
	AEnvCleanManager := &AEnvCleanManager{
		cleaner: cleaner,

		interval: duration,
		ctx:      ctx,
		cancel:   cancel,
	}
	return AEnvCleanManager
}

// Start starts the cleanup service
func (cm *AEnvCleanManager) Start() {
	log.Printf("Starting cleanup service with interval: %v", cm.interval)
	// Execute cleanup immediately
	cm.cleaner.cleanup()

	// Start timer
	ticker := time.NewTicker(cm.interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				cm.cleaner.cleanup()
			case <-cm.ctx.Done():
				log.Println("Cleanup service stopped")
				return
			}
		}
	}()
}

// Stop stops the cleanup service
func (cm *AEnvCleanManager) Stop() {
	cm.cancel()
}

// KubeCleaner cleanup service responsible for periodically cleaning expired EnvInstances
type KubeCleaner struct {
	scheduleClient *ScheduleClient
}

// NewCleanupService
func NewKubeCleaner(scheduleClient *ScheduleClient) *KubeCleaner {
	return &KubeCleaner{
		scheduleClient: scheduleClient,
	}
}

// cleanup executes cleanup task
func (cs *KubeCleaner) cleanup() {
	log.Println("Starting cleanup task...")
	// Get all EnvInstances
	envInstances, err := cs.scheduleClient.FilterPods()
	if err != nil {
		log.Printf("Failed to get env instances: %v", err)
		return
	}
	if envInstances == nil || len(*envInstances) == 0 {
		log.Println("No env instances found")
		return
	}

	var deletedCount int

	for _, instance := range *envInstances {
		// Skip terminated instances
		if instance.Status == "Terminated" {
			continue
		}
		deleted, err := cs.scheduleClient.DeletePod(instance.ID)
		if err != nil {
			log.Printf("Failed to delete instance %s: %v", instance.ID, err)
			continue
		}
		if deleted {
			deletedCount++
			log.Printf("Successfully deleted instance %s", instance.ID)
		} else {
			log.Printf("Instance %s was not deleted (may already be deleted)", instance.ID)
		}
	}
	log.Printf("Cleanup task completed. Deleted %d expired instances", deletedCount)
}

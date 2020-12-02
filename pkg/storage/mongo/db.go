package mongo

import (
	"context"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sk000f/metrix/pkg/collector"
)

var mongoOnce sync.Once
var clientInstance *mongo.Client
var clientInstanceError error

// DB represents a MongoDB client
type DB struct {
	ConnStr string
}

// SaveProjects saves Projects into the MongoDB database
func (m *DB) SaveProjects(p []*collector.Project) {
	for _, proj := range p {
		mP := Project{
			ProjectID:         proj.ID,
			Name:              proj.Name,
			PathWithNamespace: proj.PathWithNamespace,
			Namespace:         proj.Namespace,
			WebURL:            proj.WebURL,
		}

		m.UpdateProject(mP)
	}
}

// SaveDeployment saves a Deployment into the MongoDB database
func (m *DB) SaveDeployment(d *collector.Deployment) {
	mD := Deployment{
		DeploymentID:     d.ID,
		Status:           d.Status,
		EnvironmentName:  d.EnvironmentName,
		ProjectID:        d.ProjectID,
		ProjectName:      d.ProjectName,
		ProjectNamespace: d.ProjectNamespace,
		PipelineID:       d.PipelineID,
	}
	m.UpdateDeployment(mD)
}

// Project represents metrix view of a project object
type Project struct {
	ID                primitive.ObjectID `bson:"_id"`
	ProjectID         int                `bson:"project_id"`
	Name              string             `bson:"name"`
	PathWithNamespace string             `bson:"path_with_namespace"`
	Namespace         string             `bson:"namespace"`
	WebURL            string             `bson:"web_url"`
	GroupName         string             `bson:"group_name"`
}

// Deployment represents metrix view of a deployment object
type Deployment struct {
	ID               primitive.ObjectID `bson:"_id"`
	DeploymentID     int                `bson:"deployment_id"`
	Status           string             `bson:"status"`
	EnvironmentName  string             `bson:"envrionment_name"`
	ProjectID        int                `bson:"project_id"`
	ProjectName      string             `bson:"project_name"`
	ProjectNamespace string             `bson:"project_namespace"`
	PipelineID       int                `bson:"pipeline_id"`
}

// UpdateProject adds or updates the specified project in the MongoDB database
func (m *DB) UpdateProject(p Project) {

	c, err := m.GetMongoClient()
	if err != nil {
		log.Fatal(err)
	}

	collection := c.Database("metrix").Collection("projects")

	filter := bson.M{"project_id": p.ProjectID}
	updateOpts := options.Update().SetUpsert(true)

	update := bson.M{
		"$set": bson.M{
			"project_id":          p.ProjectID,
			"name":                p.Name,
			"path_with_namespace": p.PathWithNamespace,
			"namespace":           p.Namespace,
			"web_url":             p.WebURL,
		},
	}
	_, err = collection.UpdateOne(context.TODO(), filter, update, updateOpts)
	if err != nil {
		log.Fatalf("Error updating Project: %v", err.Error())
	}
}

// UpdateDeployment adds or updates the specified deployment in the MongoDB database
func (m *DB) UpdateDeployment(d Deployment) {

	c, err := m.GetMongoClient()
	if err != nil {
		log.Fatal(err)
	}

	collection := c.Database("metrix").Collection("deployments")

	filter := bson.M{"deployment_id": d.DeploymentID}
	updateOpts := options.Update().SetUpsert(true)

	update := bson.M{
		"$set": bson.M{
			"deployment_id":     d.DeploymentID,
			"status":            d.Status,
			"environment_name":  d.EnvironmentName,
			"project_id":        d.ProjectID,
			"project_name":      d.ProjectName,
			"project_namespace": d.ProjectNamespace,
			"pipeline_id":       d.PipelineID,
		},
	}
	_, err = collection.UpdateOne(context.TODO(), filter, update, updateOpts)
	if err != nil {
		log.Fatalf("Error updating Deployment: %v", err.Error())
	}
}

// GetMongoClient creates or returns existing MongoDB client
func (m *DB) GetMongoClient() (*mongo.Client, error) {

	mongoOnce.Do(func() {

		clientOptions := options.Client().ApplyURI(m.ConnStr)

		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatal(err)
			clientInstanceError = err
		}

		clientInstance = client
	})

	return clientInstance, clientInstanceError
}

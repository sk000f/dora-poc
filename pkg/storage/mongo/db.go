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
			NameWithNamespace: proj.NameWithNamespace,
			WebURL:            proj.WebURL,
		}

		m.UpdateProject(mP)
	}
}

// SaveDeployment saves a Deployment into the MongoDB database
func (m *DB) SaveDeployment(d *collector.Deployment) {
	//m.DeploymentData = append(m.DeploymentData, d)
}

// Project represents metrix view of a project object
type Project struct {
	ID                primitive.ObjectID `bson:"_id"`
	ProjectID         int                `bson:"project_id"`
	Name              string             `bson:"name"`
	NameWithNamespace string             `bson:"name_with_namespace"`
	WebURL            string             `bson:"web_url"`
	GroupName         string             `bson:"group_name"`
}

// Deployment represents metrix view of a deployment object
type Deployment struct {
	ID              primitive.ObjectID `bson:"_id"`
	DeploymentID    int                `bson:"deployment_id"`
	Status          string             `bson:"status"`
	EnvironmentName string             `bson:"envrionment_name"`
	PipelineID      int                `bson:"pipeline_id"`
	ProjectName     string             `bson:"project_name"`
	GroupName       string             `bson:"group_name"`
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
			"name_with_namespace": p.NameWithNamespace,
			"web_url":             p.WebURL,
		},
	}
	_, err = collection.UpdateOne(context.TODO(), filter, update, updateOpts)
	if err != nil {
		log.Fatalf("Error updating Project: %v", err.Error())
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

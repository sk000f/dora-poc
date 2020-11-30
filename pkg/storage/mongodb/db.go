package mongodb

import (
	"context"
	"fmt"
	"log"
	"sync"

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

		m.InsertProject(mP)
	}
}

// SaveDeployment saves a Deployment into the MongoDB database
func (m *DB) SaveDeployment(d *collector.Deployment) {
	//m.DeploymentData = append(m.DeploymentData, d)
}

// Project represents metrix view of a project object
type Project struct {
	ID                primitive.ObjectID `json:"_id"`
	ProjectID         int                `json:"projectId"`
	Name              string             `json:"name"`
	NameWithNamespace string             `json:"nameWithNamespace"`
	WebURL            string             `json:"webURL"`
	GroupName         string             `json:"groupName"`
}

// Deployment represents metrix view of a deployment object
type Deployment struct {
	ID              primitive.ObjectID `json:"_id"`
	DeploymentID    int                `json:"deploymentId"`
	Status          string             `json:"status,omitempty"`
	EnvironmentName string             `json:"envrionmentName,omitempty"`
	PipelineID      int                `json:"pipelineID,omitempty"`
	ProjectName     string             `json:"projectName,omitempty"`
	GroupName       string             `json:"groupName,omitempty"`
}

// InsertProject adds the specified project to the MongoDB database
func (m *DB) InsertProject(p Project) {

	c, err := m.GetMongoClient()
	if err != nil {
		log.Fatal(err)
	}

	collection := c.Database("metrix").Collection("projects")

	upsertResult, err := collection.InsertOne(context.TODO(), p)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Inserted Project with ObjectID: %v\n", upsertResult.InsertedID)
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

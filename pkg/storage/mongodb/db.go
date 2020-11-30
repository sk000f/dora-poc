package mongodb

import (
	"context"
	"fmt"
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

		m.InsertProject(mP)
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
	ID              primitive.ObjectID `json:"_id"`
	DeploymentID    int                `json:"deployment_id"`
	Status          string             `json:"status"`
	EnvironmentName string             `json:"envrionment_name"`
	PipelineID      int                `json:"pipeline_id"`
	ProjectName     string             `json:"project_name"`
	GroupName       string             `json:"group_name"`
}

func (m *DB) projectExists(p Project) (bool, error) {

	c, err := m.GetMongoClient()
	if err != nil {
		log.Fatal(err)
	}

	collection := c.Database("metrix").Collection("projects")

	var eP Project
	findOpts := options.FindOne()
	filter := bson.M{"project_id": p.ProjectID}
	err = collection.FindOne(context.TODO(), filter, findOpts).Decode(&eP)

	if err == mongo.ErrNoDocuments {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// InsertProject adds the specified project to the MongoDB database
func (m *DB) InsertProject(p Project) {

	c, err := m.GetMongoClient()
	if err != nil {
		log.Fatal(err)
	}

	collection := c.Database("metrix").Collection("projects")

	// find if project exists
	// exists, err := m.projectExists(p)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	exists := true
	// if it does then update it
	if exists {

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
		updateResult, err := collection.UpdateOne(context.TODO(), filter, update, updateOpts)
		if err != nil {
			log.Fatalf("Error updating Project: %v", err.Error())
		}

		fmt.Printf("Updated Projects: %v\n", updateResult.MatchedCount)

	} else { // if it doesn't, set object id and insert it

		// p.ID = primitive.NewObjectID()

		// insertResult, err := collection.InsertOne(context.TODO(), p)
		// if err != nil {
		// 	log.Fatalf("Error inserting new Project: %v", err.Error())
		// }

		// fmt.Printf("Created Project: %v\n", insertResult.InsertedID)
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

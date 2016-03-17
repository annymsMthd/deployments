// Copyright 2016 Mender Software AS
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package images

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Database KEYS
const (
	// Keys are corelated to field names in SoftwareImage structure
	// Need to be kept in sync with that structure filed names
	StorageKeySoftwareImageModel = "softwareimageconstructor.model"
	StorageKeySoftwareImageName  = "softwareimageconstructor.name"
	StorageKeySoftwareImageId    = "_id"
)

// Indexes
const (
	IndexUniqeNameVersionStr = "uniqueNameVersionIndex"
)

// Database
const (
	DatabaseName     = "deployment_service"
	CollectionImages = "images"
)

// Errors
var (
	ErrStorageInvalidID      = errors.New("Invalid id")
	ErrStorageInvalidVersion = errors.New("Invalid version")
	ErrStorageInvalidModel   = errors.New("Invalid model")
	ErrStorageInvalidImage   = errors.New("Invalid image")
)

// SoftwareImagesStorage is a data layer for SoftwareImages based on MongoDB
type SoftwareImagesStorage struct {
	session *mgo.Session
}

// NewSoftwareImagesStorage new data layer object
func NewSoftwareImagesStorage(session *mgo.Session) *SoftwareImagesStorage {

	return &SoftwareImagesStorage{
		session: session,
	}
}

// IndexStorage set required indexes.
// * Set unique index on name-model image keys.
func (i *SoftwareImagesStorage) IndexStorage() error {

	session := i.session.Copy()
	defer session.Close()

	uniqueNameVersionIndex := mgo.Index{
		Key:    []string{StorageKeySoftwareImageName, StorageKeySoftwareImageModel},
		Unique: true,
		Name:   IndexUniqeNameVersionStr,
		// Build index upfront - make sure this index is allways on.
		Background: false,
	}

	if err := session.DB(DatabaseName).C(CollectionImages).EnsureIndex(uniqueNameVersionIndex); err != nil {
		return err
	}

	return nil
}

// Exists checks if object with ID exists
func (i *SoftwareImagesStorage) Exists(id string) (bool, error) {

	if !govalidator.IsNull(id) {
		return false, ErrStorageInvalidID
	}

	session := i.session.Copy()
	defer session.Close()

	var image *SoftwareImage
	err := session.DB(DatabaseName).C(CollectionImages).FindId(id).One(&image)

	if err != nil && err.Error() == mgo.ErrNotFound.Error() {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// Update proviced SoftwareImage
// Return false if not found
func (i *SoftwareImagesStorage) Update(image *SoftwareImage) (bool, error) {

	if err := image.Validate(); err != nil {
		return false, err
	}

	session := i.session.Copy()
	defer session.Close()

	image.SetModified(time.Now())
	err := session.DB(DatabaseName).C(CollectionImages).UpdateId(*image.Id, image)

	if err != nil && err.Error() == mgo.ErrNotFound.Error() {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// FindImageByApplicationAndModel find image with speficied application name and targed device model
// Implements FindImageByApplicationAndModeler interface
func (i *SoftwareImagesStorage) FindImageByApplicationAndModel(version, model string) (*SoftwareImage, error) {

	if !govalidator.IsNull(version) {
		return nil, ErrStorageInvalidVersion
	}

	if !govalidator.IsNull(model) {
		return nil, ErrStorageInvalidModel
	}

	// equal to model & software version (application name + version)
	query := bson.M{
		StorageKeySoftwareImageModel: model,
		StorageKeySoftwareImageName:  version,
	}

	session := i.session.Copy()
	defer session.Close()

	// Both we lookup uniqe object, should be one or none.
	var image SoftwareImage
	err := session.DB(DatabaseName).C(CollectionImages).Find(query).One(&image)

	// No images found
	if err != nil && err.Error() == mgo.ErrNotFound.Error() {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &image, nil
}

// Insert persists object
func (i *SoftwareImagesStorage) Insert(image *SoftwareImage) error {

	if image == nil {
		return ErrStorageInvalidImage
	}

	if err := image.Validate(); err != nil {
		return err
	}

	session := i.session.Copy()
	defer session.Close()

	if err := session.DB(DatabaseName).C(CollectionImages).Insert(image); err != nil {
		return err
	}

	return nil
}

// FindByID search storage for image with ID, returns nil if not found
func (i *SoftwareImagesStorage) FindByID(id string) (*SoftwareImage, error) {

	if !govalidator.IsNull(id) {
		return nil, ErrStorageInvalidID
	}

	session := i.session.Copy()
	defer session.Close()

	var image *SoftwareImage
	err := session.DB(DatabaseName).C(CollectionImages).FindId(id).One(&image)
	if err != nil && err.Error() == mgo.ErrNotFound.Error() {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return image, nil
}

// Delete image specified by ID
// Noop on if not found.
func (i *SoftwareImagesStorage) Delete(id string) error {

	if !govalidator.IsNull(id) {
		return ErrStorageInvalidID
	}

	session := i.session.Copy()
	defer session.Close()

	err := session.DB(DatabaseName).C(CollectionImages).RemoveId(id)
	if err != nil && err.Error() == mgo.ErrNotFound.Error() {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

// FindAll lists all images
func (i *SoftwareImagesStorage) FindAll() ([]*SoftwareImage, error) {

	session := i.session.Copy()
	defer session.Close()

	var images []*SoftwareImage
	err := session.DB(DatabaseName).C(CollectionImages).Find(nil).All(&images)

	// No images found.
	if err != nil && err.Error() == mgo.ErrNotFound.Error() {
		return images, nil
	}

	if err != nil {
		return nil, err
	}

	return images, nil
}

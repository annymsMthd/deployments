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
package mvc_old

type CreateModeler interface {
	NewObject() interface{}
	Validate(interface{}) error
	Create(interface{}) (string, error)
}

type GetObjectModeler interface {
	GetObject(id string) (interface{}, error)
}

type DeleteObjectModeler interface {
	DeleteObject(id string) error
}

type ListObjectsModeler interface {
	ListObjects(map[string]string) (interface{}, error)
}

type EditObjectModeler interface {
	NewObject() interface{}
	Validate(interface{}) error
	EditObject(id string, object interface{}) (bool, error)
}

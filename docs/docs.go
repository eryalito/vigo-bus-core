// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/lines": {
            "get": {
                "description": "Provide a list of all the lines",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Bus"
                ],
                "summary": "List all of the lines",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/api.Line"
                            }
                        }
                    }
                }
            }
        },
        "/api/stops": {
            "get": {
                "description": "Provide a list of all the stops",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Bus"
                ],
                "summary": "List all of the stops",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/api.Stop"
                            }
                        }
                    }
                }
            }
        },
        "/api/stops/find": {
            "get": {
                "description": "Provide a list of stops that match the text in their name",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Bus"
                ],
                "summary": "Find a stop by text in its name",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Text to search for in stop name",
                        "name": "text",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/api.Stop"
                            }
                        }
                    }
                }
            }
        },
        "/api/stops/find/location": {
            "get": {
                "description": "Provide a list of stops in a given radius around a location",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Bus"
                ],
                "summary": "Find a stop by its location",
                "parameters": [
                    {
                        "type": "number",
                        "description": "Latitude",
                        "name": "lat",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "number",
                        "description": "Longitude",
                        "name": "lon",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "number",
                        "description": "Radius in meters",
                        "name": "radius",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/api.Stop"
                            }
                        }
                    }
                }
            }
        },
        "/api/stops/find/location/image": {
            "get": {
                "description": "Provide the nearby stops for a location and return a PNG image and JSON array",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Bus"
                ],
                "summary": "Get the nearby stops as a PNG image and JSON array",
                "parameters": [
                    {
                        "type": "number",
                        "description": "Latitude",
                        "name": "lat",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "number",
                        "description": "Longitude",
                        "name": "lon",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "number",
                        "description": "Radius in meters",
                        "name": "radius",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Limit of stops to return, default 9",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.NearbyStops"
                        }
                    }
                }
            }
        },
        "/api/stops/{stop_number}": {
            "get": {
                "description": "Provide a stop by its number",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Bus"
                ],
                "summary": "Get a stop by its number",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Stop Number",
                        "name": "stop_number",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.Stop"
                        }
                    }
                }
            }
        },
        "/api/stops/{stop_number}/schedule": {
            "get": {
                "description": "Provide the schedule for a stop",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Bus"
                ],
                "summary": "Get the schedule for a stop",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Stop Number",
                        "name": "stop_number",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.StopSchedule"
                        }
                    }
                }
            }
        },
        "/api/users/{provider}/{uuid}": {
            "get": {
                "description": "Provide a user by its UUID for a specific provider",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Identity"
                ],
                "summary": "Get a user by its UUID for a specific provider",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Provider",
                        "name": "provider",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "UUID",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.Identity"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Identity"
                ],
                "summary": "Create a new user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Provider",
                        "name": "provider",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "UUID",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.Identity"
                        }
                    }
                }
            }
        },
        "/api/users/{provider}/{uuid}/favorite_stops/{stop_number}": {
            "post": {
                "description": "Add a favorite stop to a user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Identity"
                ],
                "summary": "Add a favorite stop to a user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Provider",
                        "name": "provider",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "UUID",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Stop Number",
                        "name": "stop_number",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.Identity"
                        }
                    }
                }
            },
            "delete": {
                "description": "Remove a favorite stop from a user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Identity"
                ],
                "summary": "Remove a favorite stop from a user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Provider",
                        "name": "provider",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "UUID",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Stop Number",
                        "name": "stop_number",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.Identity"
                        }
                    }
                }
            }
        },
        "/api/users/{provider}/{uuid}/metadata": {
            "put": {
                "description": "Update the metadata of a user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Identity"
                ],
                "summary": "Update the metadata of a user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Provider",
                        "name": "provider",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "UUID",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Metadata",
                        "name": "metadata",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.Identity"
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "Health endpoint",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Health"
                ],
                "summary": "Health endpoint",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.Identity": {
            "type": "object",
            "properties": {
                "favorite_stops": {
                    "description": "FavoriteStops is a list of the user's favorite stops",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/api.Stop"
                    }
                },
                "id": {
                    "description": "ID is the unique identifier of the identity",
                    "type": "integer"
                },
                "metadata": {
                    "description": "Metadata is a genric string that holds additional information about the identity",
                    "type": "string"
                },
                "provider": {
                    "description": "Provider is the type of the identity provider",
                    "allOf": [
                        {
                            "$ref": "#/definitions/api.ProviderType"
                        }
                    ]
                },
                "uuid": {
                    "description": "UUID is the unique identifier of the identity, usually provided by the auth provider",
                    "type": "string"
                }
            }
        },
        "api.Line": {
            "type": "object",
            "properties": {
                "id": {
                    "description": "ID is the unique identifier of the line",
                    "type": "integer"
                },
                "name": {
                    "description": "Name is the name of the line provided by the bus company",
                    "type": "string"
                }
            }
        },
        "api.NearbyStops": {
            "type": "object",
            "properties": {
                "image": {
                    "type": "string"
                },
                "origin": {
                    "type": "object",
                    "properties": {
                        "lat": {
                            "type": "number"
                        },
                        "lon": {
                            "type": "number"
                        }
                    }
                },
                "radius": {
                    "type": "number"
                },
                "stops": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/api.Stop"
                    }
                }
            }
        },
        "api.ProviderType": {
            "type": "string",
            "enum": [
                "telegram"
            ],
            "x-enum-varnames": [
                "ProviderTypeTelegram"
            ]
        },
        "api.Schedule": {
            "type": "object",
            "properties": {
                "line": {
                    "description": "Line is the line that the schedule is for",
                    "allOf": [
                        {
                            "$ref": "#/definitions/api.Line"
                        }
                    ]
                },
                "route": {
                    "description": "Route is the route that the schedule is for",
                    "type": "string"
                },
                "time": {
                    "description": "Time is the time of the schedule",
                    "type": "integer"
                }
            }
        },
        "api.Stop": {
            "type": "object",
            "properties": {
                "id": {
                    "description": "ID is the unique identifier of the stop",
                    "type": "integer"
                },
                "location": {
                    "description": "Location is the geographical location of the stop",
                    "type": "object",
                    "properties": {
                        "lat": {
                            "description": "Lat is the latitude of the stop",
                            "type": "number"
                        },
                        "lon": {
                            "description": "Lon is the longitude of the stop",
                            "type": "number"
                        }
                    }
                },
                "name": {
                    "description": "Name is the name of the stop",
                    "type": "string"
                },
                "stop_id": {
                    "description": "StopID is the number of the stop used internally by the bus company",
                    "type": "integer"
                },
                "stop_number": {
                    "description": "StopNumber is the number of the stop provided by the bus company",
                    "type": "integer"
                }
            }
        },
        "api.StopSchedule": {
            "type": "object",
            "properties": {
                "schedules": {
                    "description": "Schedules is a list of the schedules for the stop",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/api.Schedule"
                    }
                },
                "stop": {
                    "description": "Stop is the stop that the schedule is for",
                    "allOf": [
                        {
                            "$ref": "#/definitions/api.Stop"
                        }
                    ]
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "description": "\"Type 'Bearer' followed by a space and then your token.\"",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    },
    "security": [
        {
            "BearerAuth": []
        }
    ]
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Vigo Bus Core API",
	Description:      "This is the API for the Vigo Bus Core project.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

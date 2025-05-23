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
        "/daily-tasks": {
            "get": {
                "description": "Retrieve the daily task for the current date",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "DailyTasks"
                ],
                "summary": "Get today's daily task",
                "responses": {
                    "200": {
                        "description": "Successfully retrieved today's task",
                        "schema": {
                            "$ref": "#/definitions/entity.DailyTasks"
                        }
                    },
                    "400": {
                        "description": "No token",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    },
                    "401": {
                        "description": "Invalid or expired token",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    },
                    "403": {
                        "description": "User is not admin",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    },
                    "404": {
                        "description": "Record not found in DB",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new set of daily tasks for the current date",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "DailyTasks"
                ],
                "summary": "Create daily tasks",
                "parameters": [
                    {
                        "description": "Daily tasks data",
                        "name": "tasks",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.DBDailyTasks"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Successfully created today's task",
                        "schema": {
                            "$ref": "#/definitions/entity.DailyTasks"
                        }
                    },
                    "400": {
                        "description": "No token or Invalid data",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    },
                    "401": {
                        "description": "Invalid or expired token",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    },
                    "403": {
                        "description": "User is not admin",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    },
                    "409": {
                        "description": "Conflict: DailyTask already exists",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete daily tasks for the current date",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "DailyTasks"
                ],
                "summary": "Delete today's daily tasks",
                "responses": {
                    "200": {
                        "description": "Successfully deleted today's tasks",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    },
                    "400": {
                        "description": "No token or Invalid data",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    },
                    "401": {
                        "description": "Invalid or expired token",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    },
                    "403": {
                        "description": "User is not admin",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    },
                    "404": {
                        "description": "Record not found in DB",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    }
                }
            }
        },
        "/deatil-plan": {
            "post": {
                "description": "Handles season planning request and verifies user admin rights",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Season Planning"
                ],
                "summary": "Season detailed planning",
                "parameters": [
                    {
                        "description": "Season information",
                        "name": "season",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.DetailSeasonJson"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Successfully created season plan",
                        "schema": {
                            "$ref": "#/definitions/entity.DetailSeasonJson"
                        }
                    },
                    "400": {
                        "description": "No token or Invalid data",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    },
                    "401": {
                        "description": "Invalid or expired token",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    },
                    "403": {
                        "description": "User is not admin",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    },
                    "409": {
                        "description": "Conflict: Seasons are crossing",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/seasons": {
            "get": {
                "description": "Retrieve the list of all available seasons",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Season Planning"
                ],
                "summary": "Get all seasons",
                "responses": {
                    "200": {
                        "description": "Successfully retrieved all seasons",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entity.Season"
                            }
                        }
                    },
                    "400": {
                        "description": "No token or Invalid data",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    },
                    "401": {
                        "description": "Invalid or expired token",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    },
                    "403": {
                        "description": "User is not admin",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/settings": {
            "get": {
                "description": "Retrieve the game configuration settings for the current user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Settings"
                ],
                "summary": "Get game settings",
                "responses": {
                    "200": {
                        "description": "Successfully retrieved game settings",
                        "schema": {
                            "$ref": "#/definitions/entity.SettingsJson"
                        }
                    },
                    "400": {
                        "description": "No token or Invalid data",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    },
                    "401": {
                        "description": "Invalid or expired token",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    },
                    "403": {
                        "description": "User is not admin",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Update game configuration settings for authenticated user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Settings"
                ],
                "summary": "Update game settings",
                "parameters": [
                    {
                        "description": "Game settings object",
                        "name": "settings",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.SettingsJson"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Successfully updated settings",
                        "schema": {
                            "$ref": "#/definitions/entity.DetailSeasonJson"
                        }
                    },
                    "400": {
                        "description": "No token or Invalid data",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    },
                    "401": {
                        "description": "Invalid or expired token",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    },
                    "403": {
                        "description": "User is not admin",
                        "schema": {
                            "$ref": "#/definitions/entity.Response"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/statistic/players": {
            "get": {
                "description": "Get season statistic",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Statistic"
                ],
                "summary": "Get season statistic",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/statistic/seasons": {
            "get": {
                "description": "Get seasons statistic",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Statistic"
                ],
                "summary": "Get seasons statistic",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.DBDailyTasks": {
            "type": "object",
            "properties": {
                "game-task-reward": {
                    "type": "integer",
                    "minimum": 0,
                    "example": 10
                },
                "games-amount": {
                    "type": "integer",
                    "minimum": 0,
                    "example": 5
                },
                "referals-amount": {
                    "type": "integer",
                    "minimum": 0,
                    "example": 10
                },
                "referals-task-reward": {
                    "type": "integer",
                    "minimum": 0,
                    "example": 10
                }
            }
        },
        "entity.DailyTasks": {
            "type": "object",
            "properties": {
                "game-task-reward": {
                    "type": "integer",
                    "example": 10
                },
                "games-amount": {
                    "type": "integer",
                    "example": 5
                },
                "referals-amount": {
                    "type": "integer",
                    "example": 10
                },
                "referals-task-reward": {
                    "type": "integer",
                    "example": 10
                },
                "task-date": {
                    "type": "string",
                    "example": "2023-05-15T00:00:00Z"
                }
            }
        },
        "entity.DetailSeasonJson": {
            "type": "object",
            "required": [
                "end-date",
                "end-time",
                "start-date",
                "start-time"
            ],
            "properties": {
                "end-date": {
                    "type": "string",
                    "format": "date",
                    "example": "31-08-2034"
                },
                "end-time": {
                    "type": "string",
                    "example": "01-20-30"
                },
                "fund": {
                    "type": "integer",
                    "minimum": 0,
                    "example": 5000
                },
                "start-date": {
                    "type": "string",
                    "format": "date",
                    "example": "01-06-2024"
                },
                "start-time": {
                    "type": "string",
                    "example": "00-20-30"
                }
            }
        },
        "entity.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "entity.Response": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "entity.Season": {
            "type": "object",
            "properties": {
                "end-date": {
                    "type": "string"
                },
                "fund": {
                    "type": "integer"
                },
                "start-date": {
                    "type": "string"
                }
            }
        },
        "entity.SettingsJson": {
            "type": "object",
            "properties": {
                "lose-amount": {
                    "type": "number",
                    "minimum": 0,
                    "example": 5.97
                },
                "waiting-time": {
                    "type": "integer",
                    "minimum": 0,
                    "example": 3
                },
                "win-amount": {
                    "type": "number",
                    "minimum": 0,
                    "example": 10.05
                }
            }
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8004",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Admin service",
	Description:      "Admin server API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

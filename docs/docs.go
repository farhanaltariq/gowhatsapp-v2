// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "fiber@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login": {
            "post": {
                "security": [
                    {
                        "Authorization": []
                    }
                ],
                "description": "Login",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Login",
                "parameters": [
                    {
                        "description": "Login Data",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Login"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.AuthenticationResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/common.ResponseMessage"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "Register a new user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Register",
                "parameters": [
                    {
                        "description": "Register data",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Register"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/common.ResponseMessage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/common.ResponseMessage"
                        }
                    }
                }
            }
        },
        "/message": {
            "post": {
                "description": "Send New Message to Desired Number",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Message"
                ],
                "summary": "Send Message",
                "parameters": [
                    {
                        "description": "Message data",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Message"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/common.ResponseMessage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/common.ResponseMessage"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "common.ResponseMessage": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "error": {
                    "type": "boolean"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "models.AuthenticationResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "expired_at": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "models.Login": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "usernameOrEmail": {
                    "type": "string"
                }
            }
        },
        "models.Message": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "number": {
                    "type": "integer"
                }
            }
        },
        "models.Register": {
            "type": "object",
            "properties": {
                "confirmPassword": {
                    "type": "string"
                },
                "country": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "Authorization": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:3000",
	BasePath:         "/api",
	Schemes:          []string{"http", "https"},
	Title:            "Fiber Boilerplate API",
	Description:      "This is a sample swagger for Fiber",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

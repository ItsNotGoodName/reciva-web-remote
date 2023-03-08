// Code generated by swaggo/swag. DO NOT EDIT
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
        "/build": {
            "get": {
                "tags": [
                    "build"
                ],
                "summary": "Get build",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Build"
                        }
                    }
                }
            }
        },
        "/presets": {
            "get": {
                "tags": [
                    "presets"
                ],
                "summary": "List presets",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Preset"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.HTTPError"
                        }
                    }
                }
            },
            "post": {
                "tags": [
                    "presets"
                ],
                "summary": "Update preset",
                "parameters": [
                    {
                        "description": "Preset",
                        "name": "preset",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Preset"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/http.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.HTTPError"
                        }
                    }
                }
            }
        },
        "/presets/{url}": {
            "get": {
                "tags": [
                    "presets"
                ],
                "summary": "Get preset",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Preset URL",
                        "name": "url",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Preset"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/http.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.HTTPError"
                        }
                    }
                }
            }
        },
        "/radios": {
            "get": {
                "tags": [
                    "radios"
                ],
                "summary": "List radios",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            },
            "post": {
                "tags": [
                    "radios"
                ],
                "summary": "Discover radios",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "409": {
                        "description": "Discovery already in progress",
                        "schema": {
                            "$ref": "#/definitions/http.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.HTTPError"
                        }
                    }
                }
            }
        },
        "/radios/{uuid}": {
            "get": {
                "tags": [
                    "radios"
                ],
                "summary": "Get radio",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Radio UUID",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/radios/{uuid}/subscription": {
            "post": {
                "tags": [
                    "radios"
                ],
                "summary": "Refresh radio subscription",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Radio UUID",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/radios/{uuid}/volume": {
            "post": {
                "tags": [
                    "radios"
                ],
                "summary": "Refresh radio volume",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Radio UUID",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/states": {
            "get": {
                "tags": [
                    "states"
                ],
                "summary": "List states",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/state.State"
                            }
                        }
                    }
                }
            }
        },
        "/states/{uuid}": {
            "get": {
                "tags": [
                    "states"
                ],
                "summary": "Get state",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Radio UUID",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/state.State"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.HTTPError"
                        }
                    }
                }
            },
            "patch": {
                "tags": [
                    "states"
                ],
                "summary": "Patch state",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Radio UUID",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Patch state",
                        "name": "state",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/http.PatchState"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/state.State"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.HTTPError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "http.HTTPError": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "http.PatchState": {
            "type": "object",
            "properties": {
                "audio_source": {
                    "type": "string"
                },
                "power": {
                    "type": "boolean"
                },
                "preset": {
                    "type": "integer"
                },
                "volume": {
                    "type": "integer"
                }
            }
        },
        "model.Build": {
            "type": "object",
            "properties": {
                "built_by": {
                    "type": "string"
                },
                "commit": {
                    "type": "string"
                },
                "date": {
                    "type": "string"
                },
                "release_url": {
                    "type": "string"
                },
                "summary": {
                    "type": "string"
                },
                "version": {
                    "type": "string"
                }
            }
        },
        "model.Preset": {
            "type": "object",
            "properties": {
                "title_new": {
                    "description": "TitleNew is the overridden title.",
                    "type": "string"
                },
                "url": {
                    "description": "URL of the preset.",
                    "type": "string"
                },
                "url_new": {
                    "description": "URLNew is the overridden URL.",
                    "type": "string"
                }
            }
        },
        "state.Preset": {
            "type": "object",
            "properties": {
                "number": {
                    "description": "Number is the preset number.",
                    "type": "integer"
                },
                "title": {
                    "description": "Title of the preset.",
                    "type": "string"
                },
                "title_new": {
                    "description": "TitleNew is the overridden title.",
                    "type": "string"
                },
                "url": {
                    "description": "URL of the preset.",
                    "type": "string"
                },
                "url_new": {
                    "description": "URLNew is the overridden URL.",
                    "type": "string"
                }
            }
        },
        "state.State": {
            "type": "object",
            "properties": {
                "audio_source": {
                    "description": "AudioSource is the audio source.",
                    "type": "string"
                },
                "audio_sources": {
                    "description": "AudioSources is the list of available audio sources.",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "is_muted": {
                    "description": "IsMuted represents if the radio is muted.",
                    "type": "boolean"
                },
                "metadata": {
                    "description": "Metadata of the current playing stream.",
                    "type": "string"
                },
                "model_name": {
                    "description": "ModelName is the model name of the device.",
                    "type": "string"
                },
                "model_number": {
                    "description": "ModelNumber is the model number of the device.",
                    "type": "string"
                },
                "name": {
                    "description": "Name of the radio.",
                    "type": "string"
                },
                "power": {
                    "description": "Power represents if the radio is not in standby.",
                    "type": "boolean"
                },
                "preset_number": {
                    "description": "PresetNumber is the current preset that is playing.",
                    "type": "integer"
                },
                "presets": {
                    "description": "Presets of the radio.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/state.Preset"
                    }
                },
                "status": {
                    "description": "Status is either playing, connecting, or stopped.",
                    "type": "string"
                },
                "title": {
                    "description": "Title of the current playing stream.",
                    "type": "string"
                },
                "title_new": {
                    "description": "TitleNew is the overridden title.",
                    "type": "string"
                },
                "url": {
                    "description": "URL of the stream that is currently selected.",
                    "type": "string"
                },
                "url_new": {
                    "description": "URLNew is the overridden URL.",
                    "type": "string"
                },
                "uuid": {
                    "description": "UUID of the radio.",
                    "type": "string"
                },
                "volume": {
                    "description": "Volume of the radio.",
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Reciva Web Remote",
	Description:      "Control your legacy Reciva based internet radios (Crane, Grace Digital, Tangent, etc.) via web browser or REST API.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

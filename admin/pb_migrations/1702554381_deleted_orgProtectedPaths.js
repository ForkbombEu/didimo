// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId("2y35btjamfzcr9o");

  return dao.deleteCollection(collection);
}, (db) => {
  const collection = new Collection({
    "id": "2y35btjamfzcr9o",
    "created": "2023-10-06 07:05:25.903Z",
    "updated": "2023-12-13 11:29:08.764Z",
    "name": "orgProtectedPaths",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "bruk6bzo",
        "name": "pathRegex",
        "type": "text",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "min": null,
          "max": null,
          "pattern": ""
        }
      },
      {
        "system": false,
        "id": "4q4vbdck",
        "name": "roles",
        "type": "relation",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "collectionId": "pgsh9x4x20kdgjd",
          "cascadeDelete": false,
          "minSelect": null,
          "maxSelect": null,
          "displayFields": null
        }
      }
    ],
    "indexes": [],
    "listRule": null,
    "viewRule": null,
    "createRule": null,
    "updateRule": null,
    "deleteRule": null,
    "options": {}
  });

  return Dao(db).saveCollection(collection);
})

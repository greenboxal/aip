/* eslint-disable */
import * as types from './graphql';
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';

/**
 * Map of all GraphQL operations in the project.
 *
 * This map has several performance disadvantages:
 * 1. It is not tree-shakeable, so it will include all operations in the project.
 * 2. It is not minifiable, so the string of a GraphQL query will be multiple times inside the bundle.
 * 3. It does not support dead code elimination, so it will add unused operations.
 *
 * Therefore it is highly recommended to use the babel or swc plugin for production.
 */
const documents = {
    "\n                        subscription Sub($resourceType: String!) {\n                            resourceChanged(resourceType: $resourceType) {\n                                type\n\n                                payload {\n                                    ids\n                                }\n                            }\n                        }\n                    ": types.SubDocument,
    "\n    query GetPage($id: String!) {\n        Page(id: $id) {\n            metadata {\n                id\n            }\n            spec {\n                title\n                language\n                voice\n            }\n            status {\n                markdown\n            }\n        }\n    }\n": types.GetPageDocument,
};

/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 *
 *
 * @example
 * ```ts
 * const query = gql(`query GetUser($id: ID!) { user(id: $id) { name } }`);
 * ```
 *
 * The query argument is unknown!
 * Please regenerate the types.
 */
export function gql(source: string): unknown;

/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n                        subscription Sub($resourceType: String!) {\n                            resourceChanged(resourceType: $resourceType) {\n                                type\n\n                                payload {\n                                    ids\n                                }\n                            }\n                        }\n                    "): (typeof documents)["\n                        subscription Sub($resourceType: String!) {\n                            resourceChanged(resourceType: $resourceType) {\n                                type\n\n                                payload {\n                                    ids\n                                }\n                            }\n                        }\n                    "];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n    query GetPage($id: String!) {\n        Page(id: $id) {\n            metadata {\n                id\n            }\n            spec {\n                title\n                language\n                voice\n            }\n            status {\n                markdown\n            }\n        }\n    }\n"): (typeof documents)["\n    query GetPage($id: String!) {\n        Page(id: $id) {\n            metadata {\n                id\n            }\n            spec {\n                title\n                language\n                voice\n            }\n            status {\n                markdown\n            }\n        }\n    }\n"];

export function gql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;
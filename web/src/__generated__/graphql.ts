/* eslint-disable */
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
};

export type ProfileAptitude = {
  __typename?: 'ProfileAptitude';
  Description?: Maybe<Scalars['String']>;
};

export type RootQuery = {
  __typename?: 'RootQuery';
  agent?: Maybe<Agent>;
  agentList?: Maybe<Array<Maybe<Agent>>>;
  empty?: Maybe<Scalars['String']>;
  memory?: Maybe<Memory>;
  memoryList?: Maybe<Array<Maybe<Memory>>>;
  memory_data?: Maybe<Memory_Data>;
  memory_dataList?: Maybe<Array<Maybe<Memory_Data>>>;
  memory_segment?: Maybe<Memory_Segment>;
  memory_segmentList?: Maybe<Array<Maybe<Memory_Segment>>>;
  pipeline?: Maybe<Pipeline>;
  pipelineList?: Maybe<Array<Maybe<Pipeline>>>;
  port?: Maybe<Port>;
  portList?: Maybe<Array<Maybe<Port>>>;
  port_binding?: Maybe<Port_Binding>;
  port_bindingList?: Maybe<Array<Maybe<Port_Binding>>>;
  profile?: Maybe<Profile>;
  profileList?: Maybe<Array<Maybe<Profile>>>;
  stage?: Maybe<Stage>;
  stageList?: Maybe<Array<Maybe<Stage>>>;
  task?: Maybe<Task>;
  taskList?: Maybe<Array<Maybe<Task>>>;
  team?: Maybe<Team>;
  teamList?: Maybe<Array<Maybe<Team>>>;
};


export type RootQueryAgentArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryMemoryArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryMemory_DataArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryMemory_SegmentArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryPipelineArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryPortArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryPort_BindingArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryProfileArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryStageArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryTaskArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryTeamArgs = {
  id?: InputMaybe<Scalars['String']>;
};

export type TaskPhaseStatus = {
  __typename?: 'TaskPhaseStatus';
  StageID?: Maybe<Scalars['String']>;
  State?: Maybe<Scalars['String']>;
};

export type Agent = {
  __typename?: 'agent';
  metadata?: Maybe<GithubcomgreenboxalaippkgfordforddbResourceMetadataAgentIdAgent>;
  spec?: Maybe<GithubcomgreenboxalaippkgcollectiveAgentSpec>;
  status?: Maybe<GithubcomgreenboxalaippkgcollectiveAgentStatus>;
};

export type GithubcomgreenboxalaippkgcollectiveAgentSpec = {
  __typename?: 'githubcomgreenboxalaippkgcollectiveAgentSpec';
  extra_args?: Maybe<Array<Maybe<Scalars['String']>>>;
  given_name?: Maybe<Scalars['String']>;
  port_id?: Maybe<Scalars['String']>;
  profile_id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaippkgcollectiveAgentStatus = {
  __typename?: 'githubcomgreenboxalaippkgcollectiveAgentStatus';
  last_error?: Maybe<Scalars['String']>;
  state?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaippkgcollectivePipelineSpec = {
  __typename?: 'githubcomgreenboxalaippkgcollectivePipelineSpec';
  stages?: Maybe<Array<Maybe<Stage>>>;
};

export type GithubcomgreenboxalaippkgcollectivePortBindingSpec = {
  __typename?: 'githubcomgreenboxalaippkgcollectivePortBindingSpec';
  agent_id?: Maybe<Scalars['String']>;
  port_id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaippkgcollectivePortBindingStatus = {
  __typename?: 'githubcomgreenboxalaippkgcollectivePortBindingStatus';
  last_error?: Maybe<Scalars['String']>;
  state?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaippkgcollectivePortSpec = {
  __typename?: 'githubcomgreenboxalaippkgcollectivePortSpec';
  empty?: Maybe<Scalars['Boolean']>;
};

export type GithubcomgreenboxalaippkgcollectivePortStatus = {
  __typename?: 'githubcomgreenboxalaippkgcollectivePortStatus';
  last_error?: Maybe<Scalars['String']>;
  state?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaippkgcollectiveProfileSpec = {
  __typename?: 'githubcomgreenboxalaippkgcollectiveProfileSpec';
  directive?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaippkgcollectiveTaskSpec = {
  __typename?: 'githubcomgreenboxalaippkgcollectiveTaskSpec';
  description?: Maybe<Scalars['String']>;
  output_stage_id?: Maybe<Scalars['String']>;
  pipeline_id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaippkgcollectiveTaskStatus = {
  __typename?: 'githubcomgreenboxalaippkgcollectiveTaskStatus';
  phase?: Maybe<Scalars['String']>;
  phases?: Maybe<Array<Maybe<TaskPhaseStatus>>>;
  state?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaippkgcollectiveTeamSpec = {
  __typename?: 'githubcomgreenboxalaippkgcollectiveTeamSpec';
  manager?: Maybe<Scalars['String']>;
  members?: Maybe<Array<Maybe<Scalars['String']>>>;
};

export type GithubcomgreenboxalaippkgfordforddbResourceMetadataAgentIdAgent = {
  __typename?: 'githubcomgreenboxalaippkgfordforddbResourceMetadataAgentIDAgent';
  id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaippkgfordforddbResourceMetadataMemoryDataIdMemoryData = {
  __typename?: 'githubcomgreenboxalaippkgfordforddbResourceMetadataMemoryDataIDMemoryData';
  id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaippkgfordforddbResourceMetadataMemoryIdMemory = {
  __typename?: 'githubcomgreenboxalaippkgfordforddbResourceMetadataMemoryIDMemory';
  id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaippkgfordforddbResourceMetadataMemorySegmentIdMemorySegment = {
  __typename?: 'githubcomgreenboxalaippkgfordforddbResourceMetadataMemorySegmentIDMemorySegment';
  id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaippkgfordforddbResourceMetadataPipelineIdPipeline = {
  __typename?: 'githubcomgreenboxalaippkgfordforddbResourceMetadataPipelineIDPipeline';
  id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaippkgfordforddbResourceMetadataPortBindingIdPortBinding = {
  __typename?: 'githubcomgreenboxalaippkgfordforddbResourceMetadataPortBindingIDPortBinding';
  id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaippkgfordforddbResourceMetadataPortIdPort = {
  __typename?: 'githubcomgreenboxalaippkgfordforddbResourceMetadataPortIDPort';
  id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaippkgfordforddbResourceMetadataProfileIdProfile = {
  __typename?: 'githubcomgreenboxalaippkgfordforddbResourceMetadataProfileIDProfile';
  id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaippkgfordforddbResourceMetadataStageIdStage = {
  __typename?: 'githubcomgreenboxalaippkgfordforddbResourceMetadataStageIDStage';
  id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaippkgfordforddbResourceMetadataTaskIdTask = {
  __typename?: 'githubcomgreenboxalaippkgfordforddbResourceMetadataTaskIDTask';
  id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaippkgfordforddbResourceMetadataTeamIdTeam = {
  __typename?: 'githubcomgreenboxalaippkgfordforddbResourceMetadataTeamIDTeam';
  id?: Maybe<Scalars['String']>;
};

export type Memory = {
  __typename?: 'memory';
  branch_memory_id?: Maybe<Scalars['String']>;
  clock?: Maybe<Scalars['Int']>;
  data?: Maybe<Memory_Data>;
  height?: Maybe<Scalars['Int']>;
  metadata?: Maybe<GithubcomgreenboxalaippkgfordforddbResourceMetadataMemoryIdMemory>;
  parent_memory_id?: Maybe<Scalars['String']>;
  root_memory_id?: Maybe<Scalars['String']>;
};

export type Memory_Data = {
  __typename?: 'memory_data';
  cid?: Maybe<Scalars['String']>;
  data?: Maybe<Scalars['String']>;
  metadata?: Maybe<GithubcomgreenboxalaippkgfordforddbResourceMetadataMemoryDataIdMemoryData>;
};

export type Memory_Segment = {
  __typename?: 'memory_segment';
  metadata?: Maybe<GithubcomgreenboxalaippkgfordforddbResourceMetadataMemorySegmentIdMemorySegment>;
};

export type Pipeline = {
  __typename?: 'pipeline';
  metadata?: Maybe<GithubcomgreenboxalaippkgfordforddbResourceMetadataPipelineIdPipeline>;
  spec?: Maybe<GithubcomgreenboxalaippkgcollectivePipelineSpec>;
};

export type Port = {
  __typename?: 'port';
  metadata?: Maybe<GithubcomgreenboxalaippkgfordforddbResourceMetadataPortIdPort>;
  spec?: Maybe<GithubcomgreenboxalaippkgcollectivePortSpec>;
  status?: Maybe<GithubcomgreenboxalaippkgcollectivePortStatus>;
};

export type Port_Binding = {
  __typename?: 'port_binding';
  metadata?: Maybe<GithubcomgreenboxalaippkgfordforddbResourceMetadataPortBindingIdPortBinding>;
  spec?: Maybe<GithubcomgreenboxalaippkgcollectivePortBindingSpec>;
  status?: Maybe<GithubcomgreenboxalaippkgcollectivePortBindingStatus>;
};

export type Profile = {
  __typename?: 'profile';
  aptitudes?: Maybe<Array<Maybe<ProfileAptitude>>>;
  metadata?: Maybe<GithubcomgreenboxalaippkgfordforddbResourceMetadataProfileIdProfile>;
  spec?: Maybe<GithubcomgreenboxalaippkgcollectiveProfileSpec>;
};

export type Stage = {
  __typename?: 'stage';
  assigned_team?: Maybe<Scalars['String']>;
  depends_on?: Maybe<Array<Maybe<Scalars['String']>>>;
  id?: Maybe<Scalars['String']>;
  metadata?: Maybe<GithubcomgreenboxalaippkgfordforddbResourceMetadataStageIdStage>;
};

export type Task = {
  __typename?: 'task';
  metadata?: Maybe<GithubcomgreenboxalaippkgfordforddbResourceMetadataTaskIdTask>;
  spec?: Maybe<GithubcomgreenboxalaippkgcollectiveTaskSpec>;
  status?: Maybe<GithubcomgreenboxalaippkgcollectiveTaskStatus>;
};

export type Team = {
  __typename?: 'team';
  metadata?: Maybe<GithubcomgreenboxalaippkgfordforddbResourceMetadataTeamIdTeam>;
  spec?: Maybe<GithubcomgreenboxalaippkgcollectiveTeamSpec>;
};

export type GetAgentsQueryVariables = Exact<{ [key: string]: never; }>;


export type GetAgentsQuery = { __typename?: 'RootQuery', agentList?: Array<{ __typename?: 'agent', metadata?: { __typename?: 'githubcomgreenboxalaippkgfordforddbResourceMetadataAgentIDAgent', id?: string | null } | null } | null> | null };


export const GetAgentsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetAgents"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"agentList"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"metadata"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}}]}}]}}]} as unknown as DocumentNode<GetAgentsQuery, GetAgentsQueryVariables>;
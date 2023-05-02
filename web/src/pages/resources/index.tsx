import { useQuery } from '@apollo/client';
import { gql } from '../../__generated__';

const GET_AGENTS = gql(`
    query GetAgents {
        agentList {
            metadata {
                id
            }
        }
    }
`)


export const ResourcesPage = () => {
    const { loading, error, data } = useQuery(GET_AGENTS);

    if (loading) return <p>Loading...</p>;
    if (error) return <p>Error : {error.message}</p>;

    const resources = data.agentList.map((agent) => (
        <div key={agent.metadata.id}>
            <p>{JSON.stringify(agent)}</p>
        </div>
    ))

    return <div>
        <h1>Resources</h1>
        {resources}
    </div>
}

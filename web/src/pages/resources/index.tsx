import { useQuery, gql } from '@apollo/client';

interface Agent {
    id: string;
}

const GET_AGENTS = gql`
    query GetAgents {
        agentList {
            id
        }
    }
`

const ResourcesPage = () => {
    const { loading, error, data } = useQuery(GET_AGENTS);

    if (loading) return <p>Loading...</p>;
    if (error) return <p>Error : {error.message}</p>;

    const resources = data.agentList.map((agent: any) => {

    })

    return <div>
        <h1>Resources</h1>
    </div>
}

export {}

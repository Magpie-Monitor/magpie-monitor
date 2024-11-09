import React from 'react';
import EntriesSelector from 'components/EntriesSelector/EntriesSelector';
import { AccuracyLevel, ManagmentServiceApiInstance } from 'api/managment-service.ts';
import LinkComponent from 'components/LinkComponent/LinkComponent.tsx';
import { NodeDataRow } from 'pages/Report/NodesSection/NodesSection.tsx';

interface NodesEntriesSelectorProps {
    selectedNodes: NodeDataRow[];
    setSelectedNodes: React.Dispatch<React.SetStateAction<NodeDataRow[]>>;
    nodesToExclude: NodeDataRow[];
    onAdd: () => void;
    onClose: () => void;
    clusterId: string;
    defaultAccuracy: AccuracyLevel;
}

const NodesEntriesSelector: React.FC<NodesEntriesSelectorProps> = ({
                                                                       selectedNodes,
                                                                       setSelectedNodes,
                                                                       nodesToExclude,
                                                                       onAdd,
                                                                       onClose,
                                                                       clusterId,
                                                                       defaultAccuracy,
                                                                   }) => {
    const fetchNodes = async () => {
        try {
            const data = await ManagmentServiceApiInstance.getNodes(clusterId);
            return data.map((node) => ({
                name: node.name,
                running: node.running,
                accuracy: defaultAccuracy,
                customPrompt: '',
            }));
        } catch (error) {
            console.error('Failed to fetch nodes:', error);
            return [];
        }
    };

    const getUniqueKey = (node: NodeDataRow) => node.name;

    const columns = [
        {
            header: 'Name',
            columnKey: 'name',
            customComponent: (row: NodeDataRow) => (
                <LinkComponent to="" isRunning={row.running}>
                    {row.name}
                </LinkComponent>
            ),
        },
    ];

    return (
        <EntriesSelector<NodeDataRow>
            selectedItems={selectedNodes}
            setSelectedItems={setSelectedNodes}
            itemsToExclude={nodesToExclude}
            onAdd={onAdd}
            onClose={onClose}
            fetchData={fetchNodes}
            columns={columns}
            getUniqueKey={getUniqueKey}
            entityLabel="node"
            noEntriesMessage={<p>There is no node to add.</p>}
        />
    );
};

export default NodesEntriesSelector;

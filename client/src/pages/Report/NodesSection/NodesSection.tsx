import SectionComponent from 'components/SectionComponent/SectionComponent.tsx';
import Table, { TableColumn } from 'components/Table/Table.tsx';
import React, { useState } from 'react';
import TagButton from 'components/TagButton/TagButton.tsx';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';
import ActionButton, { ActionButtonColor } from 'components/ActionButton/ActionButton.tsx';
import OverlayComponent from 'components/OverlayComponent/OverlayComponent.tsx';
import LinkComponent from 'components/LinkComponent/LinkComponent.tsx';
import { AccuracyLevel } from 'api/managment-service';
import NodesEntriesSelector
    from 'components/EntriesSelector/NodesEntriesSelector/NodesEntriesSelector.tsx';
import CustomTag from 'components/CustomTag/CustomTag.tsx';
import CustomPromptPopup from 'components/CustomPromptPopup/CustomPromptPopup.tsx';
import KindTag from 'components/KindTag/KindTag.tsx';

export interface NodeDataRow {
    name: string;
    running: boolean;
    accuracy: AccuracyLevel;
    customPrompt: string;
    [key: string]: string | boolean | AccuracyLevel;
}

interface NodesSectionProps {
    nodes: NodeDataRow[];
    setNodes: (nodes: NodeDataRow[]) => void;
    clusterId: string;
    defaultAccuracy: AccuracyLevel;
}

const NodesSection: React.FC<NodesSectionProps> = ({
                                                       setNodes,
                                                       clusterId,
                                                       defaultAccuracy,
                                                       nodes,
                                                   }) => {
    const [showModal, setShowModal] = useState(false);
    const [selectedNodes, setSelectedNodes] = useState<NodeDataRow[]>([]);
    const [showCustomPromptPopup, setShowCustomPromptPopup] = useState(false);
    const [selectedNode, setSelectedNode] = useState<NodeDataRow | null>(null);

    const handleAddNodes = () => {
        setNodes([...nodes, ...selectedNodes]);
        setSelectedNodes([]);
        setShowModal(false);
    };

    const handleAccuracyChange = (name: string, accuracy: AccuracyLevel) => {
        setNodes(
            nodes.map((node) =>
                node.name === name ? { ...node, accuracy } : node
            )
        );
    };

    const handleCustomPromptSave = (newPrompt: string) => {
        if (selectedNode) {
            setNodes(
                nodes.map((node) =>
                    node.name === selectedNode.name
                        ? { ...node, customPrompt: newPrompt }
                        : node
                )
            );
            setShowCustomPromptPopup(false);
            setSelectedNode(null);
        }
    };

    const handleCustomPromptClick = (node: NodeDataRow) => {
        setSelectedNode(node);
        setShowCustomPromptPopup(true);
    };

    const handleDelete = (name: string) => {
        setNodes(nodes.filter((node) => node.name !== name));
    };

    const columns: Array<TableColumn<NodeDataRow>> = [
        {
            header: 'Name',
            columnKey: 'name',
            customComponent: (node: NodeDataRow) => (
                <LinkComponent to="" isRunning={node.running}>
                    {node.name}
                </LinkComponent>
            ),
        },
        {
            header: 'Accuracy',
            columnKey: 'accuracy',
            customComponent: (node: NodeDataRow) => (
                <TagButton
                    listItems={['HIGH', 'MEDIUM', 'LOW']}
                    chosenItem={node.accuracy}
                    onSelect={(item) => handleAccuracyChange(node.name, item)}
                />
            ),
        },
        {
            header: 'Custom prompt',
            columnKey: 'customPrompt',
            customComponent: (node: NodeDataRow) => (
                <CustomTag
                    name={node.customPrompt || 'Enter custom prompt...'}
                    onClick={() => handleCustomPromptClick(node)}
                />
            ),
        },
        {
            header: '',
            columnKey: '',
            customComponent: () => (
                <KindTag/>
            ),
        },
        {
            header: 'Actions',
            columnKey: 'actions',
            customComponent: (node: NodeDataRow) => (
                <ActionButton
                    onClick={() => handleDelete(node.name)}
                    description="Delete"
                    color={ActionButtonColor.RED}
                />
            ),
        },
    ];

    return (
        <SectionComponent
            icon={<SVGIcon iconName="application-icon" />}
            title={'Nodes'}
            callback={() => setShowModal(true)}
        >
            <OverlayComponent isDisplayed={showModal} onClose={() => setShowModal(false)}>
                <NodesEntriesSelector
                    selectedNodes={selectedNodes}
                    setSelectedNodes={setSelectedNodes}
                    nodesToExclude={nodes}
                    onAdd={handleAddNodes}
                    onClose={() => setShowModal(false)}
                    clusterId={clusterId}
                    defaultAccuracy={defaultAccuracy}
                />
            </OverlayComponent>

            {nodes.length === 0 ? (
                <p>No Nodes selected, please add new</p>
            ) : (
                <Table columns={columns} rows={nodes}/>
            )}

            {selectedNode && (
                <CustomPromptPopup
                    initialValue={selectedNode.customPrompt}
                    isDisplayed={showCustomPromptPopup}
                    onSave={handleCustomPromptSave}
                    onClose={() => setShowCustomPromptPopup(false)}
                />
            )}
        </SectionComponent>
    );
};

export default NodesSection;

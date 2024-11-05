import './NodesEntriesSelector.scss';
import React, {useEffect, useState} from 'react';
import Table, {TableColumn} from 'components/Table/Table';
import Checkbox from 'components/Checkbox/Checkbox';
import {ManagmentServiceApiInstance} from 'api/managment-service';
import LinkComponent from 'components/LinkComponent/LinkComponent';
import {NodeEntry} from 'pages/Report/NodesSection/NodesSection';
import ActionButton, {ActionButtonColor} from 'components/ActionButton/ActionButton.tsx';

interface NodesEntriesSelectorProps {
    selectedNodes: NodeEntry[];
    setSelectedNodes: React.Dispatch<React.SetStateAction<NodeEntry[]>>;
    nodesToExclude: NodeEntry[];
    onAdd: () => void;
    onClose: () => void;
}

const NodesEntriesSelector: React.FC<NodesEntriesSelectorProps> = ({
                                                                       selectedNodes,
                                                                       setSelectedNodes,
                                                                       nodesToExclude,
                                                                       onAdd,
                                                                       onClose,
                                                                   }) => {
    const [nodes, setNodes] = useState<NodeEntry[]>([]);
    const [selectAll, setSelectAll] = useState<boolean>(false);

    useEffect(() => {
        const fetchNodes = async () => {
            try {
                const data = await ManagmentServiceApiInstance.getNodes();
                const rows = data.map((node) => ({
                    name: node.name,
                    running: node.running,
                    accuracy: node.accuracy,
                    customPrompt: node.customPrompt,
                    updated: node.updated,
                    added: node.added,
                }));
                setNodes(rows);
            } catch (error) {
                console.error('Failed to fetch nodes:', error);
            }
        };
        fetchNodes();
    }, []);

    const availableNodes = nodes.filter(
        (node) => !nodesToExclude.some((excluded) => excluded.name === node.name)
    );

    useEffect(() => {
        setSelectAll(
            availableNodes.length > 0 && selectedNodes.length === availableNodes.length
        );
    }, [selectedNodes, availableNodes]);

    const handleSelectAllChange = () => {
        if (selectAll) {
            setSelectedNodes([]);
        } else {
            setSelectedNodes(availableNodes);
        }
        setSelectAll(!selectAll);
    };

    const handleCheckboxChange = (node: NodeEntry) => {
        setSelectedNodes((prevSelected) => {
            const isSelected = prevSelected.some(
                (selectedNode) => selectedNode.name === node.name
            );
            return isSelected
                ? prevSelected.filter((selectedNode) => selectedNode.name !== node.name)
                : [...prevSelected, node];
        });
    };

    const columns: TableColumn<NodeEntry>[] = [
        {
            header: (
                <Checkbox
                    checked={selectAll}
                    onChange={handleSelectAllChange}
                />
            ),
            columnKey: 'checkbox',
            customComponent: (row: NodeEntry) => (
                <Checkbox
                    checked={selectedNodes.some(
                        (selectedNode) => selectedNode.name === row.name
                    )}
                    onChange={() => handleCheckboxChange(row)}
                />
            ),
        },
        {
            header: 'Name',
            columnKey: 'name',
            customComponent: (row: NodeEntry) => (
                <LinkComponent href="#" isRunning={row.running}>
                    {row.name}
                </LinkComponent>
            ),
        },
    ];

    return (
        <div className="nodes-entries">
            <Table
                columns={columns}
                rows={availableNodes.map((node) => ({
                    ...node,
                    key: node.name,
                }))}
            />
            <div className="nodes-entries__button-container">
                <ActionButton
                    onClick={onAdd}
                    description="Add"
                    color={ActionButtonColor.GREEN}
                />
                <ActionButton
                    onClick={onClose}
                    description="Close"
                    color={ActionButtonColor.RED}
                />
            </div>
        </div>
    );
};

export default NodesEntriesSelector;

import './NodesEntriesSelector.scss';
import React, {useEffect, useState} from 'react';
import Table, {TableColumn} from 'components/Table/Table';
import Checkbox from 'components/Checkbox/Checkbox';
import {AccuracyLevel, ManagmentServiceApiInstance} from 'api/managment-service';
import LinkComponent from 'components/LinkComponent/LinkComponent';
import { NodeDataRow } from 'pages/Report/NodesSection/NodesSection';
import ActionButton, {ActionButtonColor} from 'components/ActionButton/ActionButton.tsx';

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
    const [nodes, setNodes] = useState<NodeDataRow[]>([]);
    const [selectAll, setSelectAll] = useState<boolean>(false);

    useEffect(() => {
        const fetchNodes = async () => {
            try {
                const data = await ManagmentServiceApiInstance.getNodes(clusterId);
                const rows = data.map((node) => ({
                    name: node.name,
                    running: node.running,
                    accuracy: defaultAccuracy,
                    customPrompt: '',
                }));
                setNodes(rows);
            } catch (error) {
                console.error('Failed to fetch nodes:', error);
            }
        };
        fetchNodes();
    }, [clusterId]);

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

    const handleCheckboxChange = (node: NodeDataRow) => {
        setSelectedNodes((prevSelected) => {
            const isSelected = prevSelected.some(
                (selectedNode) => selectedNode.name === node.name
            );
            return isSelected
                ? prevSelected.filter((selectedNode) => selectedNode.name !== node.name)
                : [...prevSelected, node];
        });
    };

    const columns: TableColumn<NodeDataRow>[] = [
        {
            header: (
                <Checkbox
                    checked={selectAll}
                    onChange={handleSelectAllChange}
                />
            ),
            columnKey: 'checkbox',
            customComponent: (row: NodeDataRow) => (
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
            customComponent: (row: NodeDataRow) => (
                <LinkComponent to="#" isRunning={row.running}>
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
            <div className="nodes-entries__buttons">
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

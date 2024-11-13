import SectionComponent from 'components/SectionComponent/SectionComponent.tsx';
import Table, {TableColumn} from 'components/Table/Table.tsx';
import React, {useEffect, useState} from 'react';
import TagButton from 'components/TagButton/TagButton.tsx';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';
import ActionButton, {ActionButtonColor} from 'components/ActionButton/ActionButton.tsx';
import OverlayComponent from 'components/OverlayComponent/OverlayComponent.tsx';
import LinkComponent from 'components/LinkComponent/LinkComponent.tsx';
import {AccuracyLevel} from 'api/managment-service';
import NodesEntriesSelector
    from 'components/EntriesSelector/NodesEntriesSelector/NodesEntriesSelector.tsx';
import CustomTag from 'components/BrandTag/CustomTag.tsx';
import CustomPromptPopup from 'components/CustomPromptPopup/CustomPromptPopup.tsx';

export interface NodeDataRow {
    name: string;
    running: boolean;
    accuracy: AccuracyLevel;
    customPrompt: string;

    [key: string]: string | boolean | AccuracyLevel;
}

interface NodesSectionProps {
    setNodes: (nodes: NodeDataRow[]) => void;
    clusterId: string;
    defaultAccuracy: AccuracyLevel;
}

const NodesSection: React.FC<NodesSectionProps> = ({setNodes, clusterId, defaultAccuracy}) => {
    const [rows, setRows] = useState<NodeDataRow[]>([]);
    const [showModal, setShowModal] = useState(false);
    const [selectedNodes, setSelectedNodes] = useState<NodeDataRow[]>([]);
    const [showCustomPromptPopup, setShowCustomPromptPopup] = useState(false);
    const [selectedApp, setSelectedApp] = useState<NodeDataRow | null>(null);

    useEffect(() => {
        setNodes(rows);
    }, [rows, setNodes]);

    const handleAddNodes = () => {
        setRows([...rows, ...selectedNodes]);
        setSelectedNodes([]);
        setShowModal(false);
    };

    const handleCloseModal = () => {
        setShowModal(false);
    };

    const handleAccuracyChange = (name: string, accuracy: AccuracyLevel) => {
        setRows((prevRows) =>
            prevRows.map((row) =>
                row.name === name ? {...row, accuracy} : row
            )
        );
    };

    const handleCustomPromptSave = (newPrompt: string) => {
        if (selectedApp) {
            setRows((prevRows) =>
                prevRows.map((row) =>
                    (row.name === selectedApp.name ? { ...row, customPrompt: newPrompt } : row))
            );
            setShowCustomPromptPopup(false);
        }
    };

    const handleCustomPromptClick = (row: NodeDataRow) => {
        setSelectedApp(row);
        setShowCustomPromptPopup(true);
    };

    const handleDelete = (name: string) => {
        setRows((prevRows) => prevRows.filter((row) => row.name !== name));
    };

    const columns: Array<TableColumn<NodeDataRow>> = [
        {
            header: 'Name',
            columnKey: 'name',
            customComponent: (row: NodeDataRow) => (
                <LinkComponent to="" isRunning={row.running}>
                    {row.name}
                </LinkComponent>
            ),
        },
        {
            header: 'Accuracy',
            columnKey: 'accuracy',
            customComponent: (row: NodeDataRow) => (
                <TagButton
                    listItems={['HIGH', 'MEDIUM', 'LOW']}
                    chosenItem={row.accuracy}
                    onSelect={(item) =>
                        handleAccuracyChange(row.name, item)
                    }
                />
            ),
        },
        {
            header: 'Custom prompt',
            columnKey: 'customPrompt',
            customComponent: (row: NodeDataRow) => (
                <CustomTag
                    name={row.customPrompt || 'Enter custom prompt...'}
                    onClick={() => handleCustomPromptClick(row)}
                />
            ),
        },
        {
            header: 'Actions',
            columnKey: 'actions',
            customComponent: (row: NodeDataRow) => (
                <ActionButton
                    onClick={() => handleDelete(row.name)}
                    description="Delete"
                    color={ActionButtonColor.RED}
                />
            ),
        },
    ];

    return (
        <SectionComponent
            icon={<SVGIcon iconName="application-icon"/>}
            title={'Nodes'}
            callback={() => setShowModal(true)}>
            <OverlayComponent
                isDisplayed={showModal}
                onClose={handleCloseModal}
            >
                <NodesEntriesSelector
                    selectedNodes={selectedNodes}
                    setSelectedNodes={setSelectedNodes}
                    nodesToExclude={rows}
                    onAdd={handleAddNodes}
                    onClose={handleCloseModal}
                    clusterId={clusterId}
                    defaultAccuracy={defaultAccuracy}
                />
            </OverlayComponent>
            {rows.length === 0 ? (
                <p>No Nodes selected, please add new</p>
            ) : (
                <Table columns={columns} rows={rows}/>
            )}

            {selectedApp && (
                <CustomPromptPopup
                    initialValue={selectedApp.customPrompt}
                    isDisplayed={showCustomPromptPopup}
                    onSave={handleCustomPromptSave}
                    onClose={() => setShowCustomPromptPopup(false)}
                />
            )}
        </SectionComponent>
    );
};

export default NodesSection;

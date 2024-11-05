import SectionComponent from 'components/SectionComponent/SectionComponent.tsx';
import Table, { TableColumn } from 'components/Table/Table.tsx';
import {useEffect, useState} from 'react';
import TagButton from 'components/TagButton/TagButton.tsx';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';
import ActionButton, { ActionButtonColor } from 'components/ActionButton/ActionButton.tsx';
import OverlayComponent from 'components/OverlayComponent/OverlayComponent.tsx';
import LinkComponent from 'components/LinkComponent/LinkComponent.tsx';
import CustomPrompt from 'components/CustomPrompt/CustomPrompt.tsx';
import { AccuracyLevel } from 'api/managment-service';
import NodesEntriesSelector from 'components/NodesEntriesSelector/NodesEntriesSelector.tsx';

export interface NodeEntry { //change to NodeDataRow
    name: string;
    running: boolean;
    accuracy: AccuracyLevel;
    customPrompt: string;
    updated: string;
    added: string;
    [key: string]: string | boolean | AccuracyLevel;
}

interface NodesSectionProps {
    setNodes: (nodes: NodeEntry[]) => void;
}

const NodesSection: React.FC<NodesSectionProps> = ({ setNodes }) => {
    const [rows, setRows] = useState<NodeEntry[]>([]);
    const [showModal, setShowModal] = useState(false);
    const [selectedNodes, setSelectedNodes] = useState<NodeEntry[]>([]);

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
                row.name === name ? { ...row, accuracy } : row
            )
        );
    };

    const handleCustomPromptChange = (name: string, customPrompt: string) => {
        setRows((prevRows) =>
            prevRows.map((row) =>
                row.name === name ? { ...row, customPrompt } : row,
            ),
        );
    };

    const handleDelete = (name: string) => {
        setRows((prevRows) => prevRows.filter((row) => row.name !== name));
    };

    const columns: Array<TableColumn<NodeEntry>> = [
        {
            header: 'Name',
            columnKey: 'name',
            customComponent: (row: NodeEntry) => (
                <LinkComponent to="" isRunning={row.running}>
                    {row.name}
                </LinkComponent>
            ),
        },
        {
            header: 'Accuracy',
            columnKey: 'accuracy',
            customComponent: (row: NodeEntry) => (
                <TagButton
                    listItems={['HIGH', 'MEDIUM', 'LOW']}
                    chosenItem={row.accuracy}
                    onSelect={(item) =>
                        handleAccuracyChange(row.name, item as AccuracyLevel)
                    }
                />
            ),
        },
        {
            header: 'Custom prompt',
            columnKey: 'customPrompt',
            customComponent: (row: NodeEntry) => (
                <CustomPrompt
                    value={row.customPrompt}
                    onChange={(value) => handleCustomPromptChange(row.name, value)}
                    className="node-section__input"
                />
            ),
        },
        { header: 'Updated at', columnKey: 'updated' },
        { header: 'Added at', columnKey: 'added' },
        {
            header: 'Actions',
            columnKey: 'actions',
            customComponent: (row: NodeEntry) => (
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
            icon={<SVGIcon iconName="application-icon" />}
            title={'Nodes'}
            callback={() => setShowModal(true)}>
            {showModal && (
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
                    />
                </OverlayComponent>
            )}
            {rows.length === 0 ? (
                <p>No Nodes selected, please add new</p>
            ) : (
                <Table columns={columns} rows={rows} />
            )}
        </SectionComponent>
    );
};

export default NodesSection;

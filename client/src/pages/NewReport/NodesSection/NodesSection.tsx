import SectionComponent from 'components/SectionComponent/SectionComponent.tsx';
import Table, { TableColumn } from 'components/Table/Table.tsx';
import { useEffect, useState } from 'react';
import TagButton from 'components/TagButton/TagButton.tsx';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';
import ActionButton, { ActionButtonColor } from 'components/ActionButton/ActionButton.tsx';
import OverlayComponent from 'components/OverlayComponent/OverlayComponent.tsx';
import LinkComponent from 'components/LinkComponent/LinkComponent.tsx';
import CustomPrompt from 'components/CustomPrompt/CustomPrompt.tsx';
import { ManagmentServiceApiInstance } from 'api/managment-service';

export interface NodeEntry {
    id: string;
    name: string;
    accuracy: 'HIGH' | 'MEDIUM' | 'LOW';
    customPrompt: string;
    updated: string;
    added: string;
    [key: string]: string;
}

const NodesSection = () => {
    const [rows, setRows] = useState<NodeEntry[]>([]);
    const [showModal, setShowModal] = useState(false);
    const [loading, setLoading] = useState(true);

    const fetchNodes = async () => {
        try {
            setLoading(true);
            const nodesData = await ManagmentServiceApiInstance.getNodes();

            const nodeRows = nodesData.map((node): NodeEntry => ({
                id: node.id,
                name: node.name,
                accuracy: node.accuracy,
                customPrompt: node.customPrompt,
                updated: node.updated,
                added: node.added,
            }));

            setRows(nodeRows);
        } catch (e: unknown) {
            console.error('Failed to fetch nodes', e);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchNodes();
    }, []);

    const handleAddClick = () => {
        setShowModal(true);
    };

    const handleCloseModal = () => {
        setShowModal(false);
    };

    const handleAccuracyChange = (id: string, accuracy: 'HIGH' | 'MEDIUM' | 'LOW') => {
        setRows((prevRows) =>
            prevRows.map((row) =>
                row.id === id ? { ...row, accuracy } : row
            )
        );
    };

    const handleCustomPromptChange = (id: string, customPrompt: string) => {
        setRows((prevRows) =>
            prevRows.map((row) =>
                row.id === id ? { ...row, customPrompt } : row
            )
        );
    };

    const handleDelete = (id: string) => {
        setRows((prevRows) => prevRows.filter((row) => row.id !== id));
    };

    const columns: Array<TableColumn<NodeEntry>> = [
        {
            header: 'Name',
            columnKey: 'name',
            customComponent: (row: NodeEntry) => (
                <LinkComponent href="#" className="node-section__link">
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
                    onSelect={(item) => handleAccuracyChange(row.id, item as 'HIGH' | 'MEDIUM' | 'LOW')}
                />
            ),
        },
        {
            header: 'Custom prompt',
            columnKey: 'customPrompt',
            customComponent: (row: NodeEntry) => (
                <CustomPrompt
                    value={row.customPrompt}
                    onChange={(value) => handleCustomPromptChange(row.id, value)}
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
                    onClick={() => handleDelete(row.id)}
                    description="Delete"
                    color={ActionButtonColor.RED}
                />
            ),
        }
    ];

    return (
        <SectionComponent
            icon={<SVGIcon iconName="application-icon" />}
            title={'Nodes'}
            callback={handleAddClick}>
            {showModal && (
                <OverlayComponent isDisplayed={showModal} onClose={handleCloseModal}>
                    <p>No nodes here (probably Wojciech dropped all of them)</p>
                </OverlayComponent>
            )}
            {loading ? (
                <p>Loading...</p>
            ) : rows.length === 0 ? (
                <p>No Nodes selected, please add new</p>
            ) : (
                <Table columns={columns} rows={rows} />
            )}
        </SectionComponent>
    );
};

export default NodesSection;

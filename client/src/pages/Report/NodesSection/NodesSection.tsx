import SectionComponent from 'components/SectionComponent/SectionComponent.tsx';
import Table, { TableColumn } from 'components/Table/Table.tsx';
import { useEffect, useState } from 'react';
import TagButton from 'components/TagButton/TagButton.tsx';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';
import ActionButton, { ActionButtonColor } from 'components/ActionButton/ActionButton.tsx';
import OverlayComponent from 'components/OverlayComponent/OverlayComponent.tsx';
import LinkComponent from 'components/LinkComponent/LinkComponent.tsx';
import CustomPrompt from 'components/CustomPrompt/CustomPrompt.tsx';
import { ManagmentServiceApiInstance, AccuracyLevel} from 'api/managment-service';

export interface NodeEntry {
    name: string;
    running: boolean;
    accuracy: AccuracyLevel;
    customPrompt: string;
    updated: string;
    added: string;
    [key: string]: string | boolean | AccuracyLevel;
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
                name: node.name,
                running: node.running,
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
                row.name === name ? { ...row, customPrompt } : row
            )
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
                <LinkComponent href="#" isRunning={row.running}>
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
                    onSelect={(item) => handleAccuracyChange(row.name, item as AccuracyLevel)}
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

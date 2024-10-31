import './NodesSection.scss';
import SectionComponent from 'components/SectionComponent/SectionComponent.tsx';
import Table, { TableColumn } from 'components/Table/Table.tsx';
import { useState } from 'react';
import TagButton from 'components/TagButton/TagButton.tsx';
import { NodeEntry } from 'api/managment-service';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';
import ActionButton, { ActionButtonColor } from 'components/ActionButton/ActionButton.tsx';
import OverlayComponent from 'components/OverlayComponent/OverlayComponent.tsx';

const MOCK_APPLICATIONS: NodeEntry[] = [
    {
        id: '1',
        name: 'alerts-api-database',
        precision: 'high',
        customPrompt: 'ignore s3 logs...',
        updated: '07.03.2024 15:32',
        added: '07.03.2024 15:32',
    },
    {
        id: '2',
        name: 'alerts-api-backend',
        precision: 'none',
        customPrompt: 'ignore s3 logs...',
        updated: '07.03.2024 15:32',
        added: '07.03.2024 15:32',
    },
];

const NodesSection = () => {
    const [rows, setRows] = useState<NodeEntry[]>(MOCK_APPLICATIONS);
    const [showModal, setShowModal] = useState(false);

    const columns: Array<TableColumn<NodeEntry>> = [
        {
            header: 'Name',
            columnKey: 'name',
            customComponent: (row: NodeEntry) => (
                <a href="#" className="application-section__link">
                    {row.name}
                </a>
            ),
        },
        {
            header: 'Precision',
            columnKey: 'precision',
            customComponent: (row: NodeEntry) => (
                <TagButton
                    listItems={['high', 'medium', 'low', 'none']}
                    chosenItem={row.precision}
                    onSelect={(precision) => handlePrecisionChange(row.id, precision)}
                />
            ),
        },
        {
            header: 'Custom prompt',
            columnKey: 'customPrompt',
            customComponent: (row: NodeEntry) => (
                <input
                    type="text"
                    className="application-section__input"
                    value={row.customPrompt}
                    onChange={(e) => handleCustomPromptChange(row.id, e.target.value)}
                />
            ),
        },
        { header: 'Updated', columnKey: 'updated' },
        { header: 'Added', columnKey: 'added' },
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

    const handleAddClick = () => {
        setShowModal(true);
    };

    const handleCloseModal = () => {
        setShowModal(false);
    };

    const handlePrecisionChange = (id: string, precision: string) => {
        setRows((prevRows) =>
            prevRows.map((row) =>
                row.id === id ? { ...row, precision } : row
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

    return (
        <SectionComponent
            icon={<SVGIcon iconName='application-icon' />}
            title={'Nodes'}
            callback={handleAddClick}>
            {showModal && <OverlayComponent onClose={handleCloseModal} />}
            <div className="nodes-section__table">
                {rows.length === 0 ? (
                    <p>No Nodes selected, please add new</p>
                ) : (
                    <Table columns={columns} rows={rows} />
                )}
            </div>
        </SectionComponent>
    );
};

export default NodesSection;

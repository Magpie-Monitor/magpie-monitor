import './ApplicationSection.scss';
import SectionComponent from 'components/SectionComponent/SectionComponent.tsx';
import Table, { TableColumn } from 'components/Table/Table.tsx';
import { useState } from 'react';
import TagButton from 'components/TagButton/TagButton.tsx'; // Import TagButton
import { ApplicationEntry } from 'api/managment-service';

const MOCK_APPLICATIONS: ApplicationEntry[] = [
    {
        id: '1',
        name: 'alerts-api-database',
        entries: '38,393',
        precision: 'high',
        customPrompt: 'ignore s3 logs...',
        updated: '07.03.2024 15:32',
        added: '07.03.2024 15:32',
    },
    {
        id: '2',
        name: 'alerts-api-backend',
        entries: '1,234',
        precision: 'none',
        customPrompt: 'ignore s3 logs...',
        updated: '07.03.2024 15:32',
        added: '07.03.2024 15:32',
    },
];

const ApplicationSection = () => {
    const [rows, setRows] = useState<ApplicationEntry[]>(MOCK_APPLICATIONS);

    const columns: Array<TableColumn<ApplicationEntry>> = [
        {
            header: 'Name',
            columnKey: 'name',
            customComponent: (row: ApplicationEntry) => (
                <a href="#" className="application-section__link">
                    {row.name}
                </a>
            ),
        },
        { header: 'Entries', columnKey: 'entries' },
        {
            header: 'Precision',
            columnKey: 'precision',
            customComponent: (row: ApplicationEntry) => (
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
            customComponent: (row: ApplicationEntry) => (
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
    ];

    // Handlers for updating precision and custom prompt fields
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

    return (
        <SectionComponent icon={'application-icon'} title={'Applications'}>
            <div className="application-section__table">
                <Table columns={columns} rows={rows} />
            </div>
        </SectionComponent>
    );
};

export default ApplicationSection;

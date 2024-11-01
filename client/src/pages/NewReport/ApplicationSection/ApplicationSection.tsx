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

interface ApplicationDataRow {
    id: string;
    name: string;
    accuracy: 'HIGH' | 'MEDIUM' | 'LOW';
    customPrompt: string;
    updated: string;
    added: string;
    [key: string]: string;
}

const ApplicationSection = () => {
    const [rows, setRows] = useState<ApplicationDataRow[]>([]);
    const [loading, setLoading] = useState(true);
    const [showModal, setShowModal] = useState(false);

    const fetchApplications = async () => {
        setLoading(true);
        try {
            const applicationsData = await ManagmentServiceApiInstance.getApplications();

            const applicationsRows = applicationsData.map(
                (application): ApplicationDataRow => ({
                    id: application.id,
                    name: application.name,
                    accuracy: application.accuracy,
                    customPrompt: application.customPrompt,
                    updated: application.updated,
                    added: application.added,
                }),
            );

            setRows(applicationsRows);
        } catch (e: unknown) {
            console.error('Failed to fetch applications', e);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchApplications();
    }, []);

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

    const handleAddClick = () => {
        setShowModal(true);
    };

    const handleCloseModal = () => {
        setShowModal(false);
    };

    const columns: Array<TableColumn<ApplicationDataRow>> = [
        {
            header: 'Name',
            columnKey: 'name',
            customComponent: (row: ApplicationDataRow) => (
                <LinkComponent href="#" className="application-section__link">
                    {row.name}
                </LinkComponent>
            ),
        },
        {
            header: 'Accuracy',
            columnKey: 'accuracy',
            customComponent: (row: ApplicationDataRow) => (
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
            customComponent: (row: ApplicationDataRow) => (
                <CustomPrompt
                    value={row.customPrompt}
                    onChange={(value) => handleCustomPromptChange(row.id, value)}
                    className="application-section__input"
                />
            ),
        },
        { header: 'Updated at', columnKey: 'updated' },
        { header: 'Added at', columnKey: 'added' },
        {
            header: 'Actions',
            columnKey: 'actions',
            customComponent: (row: ApplicationDataRow) => (
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
            icon={<SVGIcon iconName='application-icon' />}
            title={'Applications'}
            callback={handleAddClick}>
            {showModal && (
                <OverlayComponent isDisplayed={showModal} onClose={handleCloseModal}>
                    <p>No applications here (probably Wojciech dropped all of them)</p>
                </OverlayComponent>
            )}
            {loading ? (
                <p>Loading...</p>
            ) : rows.length === 0 ? (
                <p>No Applications selected, please add new</p>
            ) : (
                <Table columns={columns} rows={rows} />
            )}
        </SectionComponent>
    );
};

export default ApplicationSection;

import SectionComponent from 'components/SectionComponent/SectionComponent.tsx';
import Table, {TableColumn} from 'components/Table/Table.tsx';
import {useEffect, useState} from 'react';
import TagButton from 'components/TagButton/TagButton.tsx';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';
import ActionButton, {
    ActionButtonColor,
} from 'components/ActionButton/ActionButton.tsx';
import OverlayComponent from 'components/OverlayComponent/OverlayComponent.tsx';
import LinkComponent from 'components/LinkComponent/LinkComponent.tsx';
import CustomPrompt from 'components/CustomPrompt/CustomPrompt.tsx';
import {AccuracyLevel} from 'api/managment-service';
import ApplicationsEntriesSelector
    from 'components/ApplicationsEntriesSelector/ApplicationsEntriesSelector.tsx';

export interface ApplicationDataRow {
    name: string;
    running: boolean;
    accuracy: AccuracyLevel;
    customPrompt: string;
    updated: string;
    added: string;

    [key: string]: string | boolean | AccuracyLevel;
}

interface ApplicationSectionProps {
    setApplications: (apps: ApplicationDataRow[]) => void;
    clusterId: string;
}

const ApplicationSection: React.FC<ApplicationSectionProps> = ({setApplications, clusterId}) => {
    const [rows, setRows] = useState<ApplicationDataRow[]>([]);
    const [showModal, setShowModal] = useState(false);
    const [selectedApplications, setSelectedApplications] = useState<ApplicationDataRow[]>([]);

    useEffect(() => {
        setApplications(rows);
    }, [rows, setApplications]);

    const handleAddApplications = () => {
        setRows([...rows, ...selectedApplications]);
        setSelectedApplications([]);
        setShowModal(false);
    };

    const handleAccuracyChange = (name: string, accuracy: AccuracyLevel) => {
        setRows((prevRows) =>
            prevRows.map((row) => (row.name === name ? {...row, accuracy} : row)),
        );
    };

    const handleCustomPromptChange = (name: string, customPrompt: string) => {
        setRows((prevRows) =>
            prevRows.map((row) =>
                row.name === name ? {...row, customPrompt} : row,
            ),
        );
    };

    const handleCloseModal = () => {
        setShowModal(false);
    };

    const handleDelete = (name: string) => {
        setRows((prevRows) => prevRows.filter((row) => row.name !== name));
    };

    const columns: Array<TableColumn<ApplicationDataRow>> = [
        {
            header: 'Name',
            columnKey: 'name',
            customComponent: (row: ApplicationDataRow) => (
                <LinkComponent to="#" isRunning={row.running}>
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
                    onSelect={(item) =>
                        handleAccuracyChange(row.name, item as AccuracyLevel)
                    }
                />
            ),
        },
        {
            header: 'Custom prompt',
            columnKey: 'customPrompt',
            customComponent: (row: ApplicationDataRow) => (
                <CustomPrompt
                    value={row.customPrompt}
                    onChange={(value) => handleCustomPromptChange(row.name, value)}
                    className="application-section__input"
                />
            ),
        },
        {header: 'Updated at', columnKey: 'updated'},
        {header: 'Added at', columnKey: 'added'},
        {
            header: 'Actions',
            columnKey: 'actions',
            customComponent: (row: ApplicationDataRow) => (
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
            title={'Applications'}
            callback={() => setShowModal(true)}>
            <OverlayComponent
                isDisplayed={showModal}
                onClose={handleCloseModal}
            >
                <ApplicationsEntriesSelector
                    selectedApplications={selectedApplications}
                    setSelectedApplications={setSelectedApplications}
                    applicationsToExclude={rows}
                    onAdd={handleAddApplications}
                    onClose={handleCloseModal}
                    clusterId={clusterId}
                />
            </OverlayComponent>
            {rows.length === 0 ? (
                <p>No Applications selected, please add new</p>
            ) : (
                <Table columns={columns} rows={rows}/>
            )}
        </SectionComponent>
    );
};

export default ApplicationSection;


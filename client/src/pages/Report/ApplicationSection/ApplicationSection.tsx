import React, { useEffect, useState } from 'react';
import SectionComponent from 'components/SectionComponent/SectionComponent.tsx';
import Table, { TableColumn } from 'components/Table/Table.tsx';
import TagButton from 'components/TagButton/TagButton.tsx';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';
import ActionButton, { ActionButtonColor } from 'components/ActionButton/ActionButton.tsx';
import OverlayComponent from 'components/OverlayComponent/OverlayComponent.tsx';
import LinkComponent from 'components/LinkComponent/LinkComponent.tsx';
import CustomPromptPopup from 'components/CustomPromptPopup/CustomPromptPopup.tsx';
import { AccuracyLevel } from 'api/managment-service';
import ApplicationsEntriesSelector
    from 'components/EntriesSelector/ApplicationsEntriesSelector/ApplicationsEntriesSelector.tsx';
import CustomTag from 'components/BrandTag/CustomTag.tsx';

export interface ApplicationDataRow {
    name: string;
    running: boolean;
    accuracy: AccuracyLevel;
    customPrompt: string;
    kind: string;
    [key: string]: string | boolean | AccuracyLevel;
}

interface ApplicationSectionProps {
    setApplications: (apps: ApplicationDataRow[]) => void;
    clusterId: string;
    defaultAccuracy: AccuracyLevel;
}

const ApplicationSection: React.FC<ApplicationSectionProps> =
    ({ setApplications, clusterId, defaultAccuracy }) => {
    const [rows, setRows] = useState<ApplicationDataRow[]>([]);
    const [showModal, setShowModal] = useState(false);
    const [showCustomPromptPopup, setShowCustomPromptPopup] = useState(false);
    const [selectedApp, setSelectedApp] = useState<ApplicationDataRow | null>(null);
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
            prevRows.map((row) => (row.name === name ? { ...row, accuracy } : row))
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

    const handleCustomPromptClick = (row: ApplicationDataRow) => {
        setSelectedApp(row);
        setShowCustomPromptPopup(true);
    };

    const handleDelete = (name: string) => {
        setRows((prevRows) => prevRows.filter((row) => row.name !== name));
    };

    const columns: Array<TableColumn<ApplicationDataRow>> = [
        {
            header: 'Name',
            columnKey: 'name',
            customComponent: (row: ApplicationDataRow) => (
                <LinkComponent to="" isRunning={row.running}>
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
                    onSelect={(item) => handleAccuracyChange(row.name, item as AccuracyLevel)}
                />
            ),
        },
        {
            header: 'Custom prompt',
            columnKey: 'customPrompt',
            customComponent: (row: ApplicationDataRow) => (
                <CustomTag
                    name={row.customPrompt || 'Enter custom prompt...'}
                    onClick={() => handleCustomPromptClick(row)}
                />
            ),
        },
        {header: 'Kind', columnKey: 'kind'},
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
            icon={<SVGIcon iconName="application-icon" />}
            title="Applications"
            callback={() => setShowModal(true)}
        >
            <OverlayComponent isDisplayed={showModal} onClose={() => setShowModal(false)}>
                <ApplicationsEntriesSelector
                    selectedApplications={selectedApplications}
                    setSelectedApplications={setSelectedApplications}
                    applicationsToExclude={rows}
                    onAdd={handleAddApplications}
                    onClose={() => setShowModal(false)}
                    clusterId={clusterId}
                    defaultAccuracy={defaultAccuracy}
                />
            </OverlayComponent>

            {rows.length === 0 ? (
                <p>No Applications selected, please add new</p>
            ) : (
                <Table columns={columns} rows={rows} />
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

export default ApplicationSection;

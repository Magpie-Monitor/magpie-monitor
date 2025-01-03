import React, { useState } from 'react';
import SectionComponent from 'components/SectionComponent/SectionComponent.tsx';
import Table, { TableColumn } from 'components/Table/Table.tsx';
import TagButton from 'components/TagButton/TagButton.tsx';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';
import OverlayComponent from 'components/OverlayComponent/OverlayComponent.tsx';
import LinkComponent from 'components/LinkComponent/LinkComponent.tsx';
import CustomPromptPopup from 'components/CustomPromptPopup/CustomPromptPopup.tsx';
import {
    AccuracyLevel,
    ManagmentServiceApiInstance,
} from 'api/managment-service';

// eslint-disable-next-line
import ApplicationsEntriesSelector from 'components/EntriesSelector/ApplicationsEntriesSelector/ApplicationsEntriesSelector.tsx';
import CustomTag from 'components/CustomTag/CustomTag.tsx';
import KindTag from 'components/KindTag/KindTag.tsx';
import DeleteIconButton from 'components/DeleteIconButton/DeleteIconButton.tsx';

export interface ApplicationDataRow {
    name: string;
    running: boolean;
    accuracy: AccuracyLevel;
    customPrompt: string;
    kind: string;

    [key: string]: string | boolean | AccuracyLevel;
}

interface ApplicationSectionProps {
    applications: ApplicationDataRow[];
    setApplications: (apps: ApplicationDataRow[]) => void;
    clusterId: string;
    defaultAccuracy: AccuracyLevel;
}

const ApplicationSection: React.FC<ApplicationSectionProps> = ({
    applications,
    setApplications,
    clusterId,
    defaultAccuracy,
}) => {
    const [showModal, setShowModal] = useState(false);
    const [showCustomPromptPopup, setShowCustomPromptPopup] = useState(false);
    const [selectedApp, setSelectedApp] = useState<ApplicationDataRow | null>(
        null,
    );
    const [selectedApplications, setSelectedApplications] = useState<
        ApplicationDataRow[]
    >([]);
    const [applicationsToAdd, setApplicationsToAdd] = useState<
        ApplicationDataRow[]
    >([]);

    const loadApplications = async () => {
        try {
            const data = await ManagmentServiceApiInstance.getApplications(clusterId);
            setApplicationsToAdd(
                data.map((app) => ({
                    name: app.name,
                    running: app.running,
                    kind: app.kind,
                    accuracy: defaultAccuracy,
                    customPrompt: '',
                })),
            );
        } catch (error) {
            console.error('Failed to load applications:', error);
        }
    };

    const handleOpenModal = async () => {
        await loadApplications();
        setShowModal(true);
    };

    const handleAddApplications = () => {
        setApplications([...applications, ...selectedApplications]);
        setSelectedApplications([]);
        setShowModal(false);
    };

    const handleAccuracyChange = (name: string, accuracy: AccuracyLevel) => {
        setApplications(
            applications.map((app) =>
                app.name === name ? { ...app, accuracy } : app,
            ),
        );
    };

    const handleCustomPromptSave = (newPrompt: string) => {
        if (selectedApp) {
            setApplications(
                applications.map((app) =>
                    app.name === selectedApp.name
                        ? { ...app, customPrompt: newPrompt }
                        : app,
                ),
            );
            setShowCustomPromptPopup(false);
            setSelectedApp(null);
        }
    };

    const handleCustomPromptClick = (app: ApplicationDataRow) => {
        setSelectedApp(app);
        setShowCustomPromptPopup(true);
    };

    const handleDelete = (name: string) => {
        setApplications(applications.filter((app) => app.name !== name));
    };

    const columns: Array<TableColumn<ApplicationDataRow>> = [
        {
            header: 'Name',
            columnKey: 'name',
            customComponent: (app: ApplicationDataRow) => (
                <LinkComponent isRunning={app.running}>{app.name}</LinkComponent>
            ),
        },
        {
            header: 'Accuracy',
            columnKey: 'accuracy',
            customComponent: (app: ApplicationDataRow) => (
                <TagButton
                    listItems={['HIGH', 'MEDIUM', 'LOW']}
                    chosenItem={app.accuracy}
                    onSelect={(item) => handleAccuracyChange(app.name, item)}
                />
            ),
        },
        {
            header: 'Custom prompt',
            columnKey: 'customPrompt',
            customComponent: (app: ApplicationDataRow) => (
                <CustomTag
                    name={app.customPrompt || 'Enter custom prompt...'}
                    onClick={() => handleCustomPromptClick(app)}
                />
            ),
        },
        {
            header: 'Kind',
            columnKey: 'kind',
            customComponent: (app: ApplicationDataRow) => (
                <KindTag name={app.kind || 'unknown'} />
            ),
        },
        {
            header: 'Actions',
            columnKey: 'actions',
            customComponent: (app: ApplicationDataRow) => (
                <DeleteIconButton onClick={() => handleDelete(app.name)} />
            ),
        },
    ];

    return (
        <SectionComponent
            icon={<SVGIcon iconName="application-incident-metadata-icon" />}
            title="Applications"
            callback={() => handleOpenModal()}
        >
            <OverlayComponent
                isDisplayed={showModal}
                onClose={() => setShowModal(false)}
            >
                <ApplicationsEntriesSelector
                    selectedApplications={selectedApplications}
                    setSelectedApplications={setSelectedApplications}
                    applicationsToExclude={applications}
                    onAdd={handleAddApplications}
                    onClose={() => setShowModal(false)}
                    availableApplications={applicationsToAdd}
                />
            </OverlayComponent>

            {applications.length === 0 ? (
                <p>No Applications selected, please add new</p>
            ) : (
                <Table columns={columns} rows={applications} />
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

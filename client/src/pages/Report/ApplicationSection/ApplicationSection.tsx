import React, {useState} from 'react';
import SectionComponent from 'components/SectionComponent/SectionComponent.tsx';
import Table, {TableColumn} from 'components/Table/Table.tsx';
import TagButton from 'components/TagButton/TagButton.tsx';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';
import ActionButton, {ActionButtonColor} from 'components/ActionButton/ActionButton.tsx';
import OverlayComponent from 'components/OverlayComponent/OverlayComponent.tsx';
import LinkComponent from 'components/LinkComponent/LinkComponent.tsx';
import CustomPromptPopup from 'components/CustomPromptPopup/CustomPromptPopup.tsx';
import {AccuracyLevel} from 'api/managment-service';
import ApplicationsEntriesSelector
    from 'components/EntriesSelector/ApplicationsEntriesSelector/ApplicationsEntriesSelector.tsx';
import CustomTag from 'components/CustomTag/CustomTag.tsx';

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
    const [selectedApp, setSelectedApp] = useState<ApplicationDataRow | null>(null);
    const [selectedApplications, setSelectedApplications] = useState<ApplicationDataRow[]>([]);

    const handleAddApplications = () => {
        setApplications([...applications, ...selectedApplications]);
        setSelectedApplications([]);
        setShowModal(false);
    };

    const handleAccuracyChange = (name: string, accuracy: AccuracyLevel) => {
        setApplications(
            applications.map((app) =>
                app.name === name ? { ...app, accuracy } : app
            )
        );
    };

    const handleCustomPromptSave = (newPrompt: string) => {
        if (selectedApp) {
            setApplications(
                applications.map((app) =>
                    app.name === selectedApp.name
                        ? { ...app, customPrompt: newPrompt }
                        : app
                )
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
                <LinkComponent to="" isRunning={app.running}>
                    {app.name}
                </LinkComponent>
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
        { header: 'Kind', columnKey: 'kind' },
        {
            header: 'Actions',
            columnKey: 'actions',
            customComponent: (app: ApplicationDataRow) => (
                <ActionButton
                    onClick={() => handleDelete(app.name)}
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
                    applicationsToExclude={applications}
                    onAdd={handleAddApplications}
                    onClose={() => setShowModal(false)}
                    clusterId={clusterId}
                    defaultAccuracy={defaultAccuracy}
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

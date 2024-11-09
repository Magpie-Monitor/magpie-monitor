import React from 'react';
import EntriesSelector from 'components/EntriesSelector/EntriesSelector';
import { ApplicationDataRow } from 'pages/Report/ApplicationSection/ApplicationSection';
import { ManagmentServiceApiInstance, AccuracyLevel } from 'api/managment-service';
import LinkComponent from 'components/LinkComponent/LinkComponent';

interface ApplicationsEntriesSelectorProps {
    selectedApplications: ApplicationDataRow[];
    setSelectedApplications: React.Dispatch<React.SetStateAction<ApplicationDataRow[]>>;
    applicationsToExclude: ApplicationDataRow[];
    onAdd: () => void;
    onClose: () => void;
    clusterId: string;
    defaultAccuracy: AccuracyLevel;
}

const ApplicationsEntriesSelector: React.FC<ApplicationsEntriesSelectorProps> = ({
                                                                    selectedApplications,
                                                                    setSelectedApplications,
                                                                    applicationsToExclude,
                                                                    onAdd,
                                                                    onClose,
                                                                    clusterId,
                                                                    defaultAccuracy,
                                                                                 }) => {
    const fetchApplications = async () => {
        const data = await ManagmentServiceApiInstance.getApplications(clusterId);
        return data.map((app) => ({
            name: app.name,
            running: app.running,
            kind: app.kind,
            accuracy: defaultAccuracy,
            customPrompt: '',
        }));
    };

    const getUniqueKey = (app: ApplicationDataRow) => app.name;

    const columns = [
        {
            header: 'Name',
            columnKey: 'name',
            customComponent: (row: ApplicationDataRow) => (
                <LinkComponent to="" isRunning={row.running}>
                    {row.name}
                </LinkComponent>
            ),
        },
        { header: 'Kind', columnKey: 'kind' },
    ];

    return (
        <EntriesSelector<ApplicationDataRow>
            selectedItems={selectedApplications}
            setSelectedItems={setSelectedApplications}
            itemsToExclude={applicationsToExclude}
            onAdd={onAdd}
            onClose={onClose}
            fetchData={fetchApplications}
            columns={columns}
            getUniqueKey={getUniqueKey}
            entityLabel="application"
            noEntriesMessage={<p>There is no application to add.</p>}
        />
    );
};

export default ApplicationsEntriesSelector;

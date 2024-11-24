import React from 'react';
import EntriesSelector from 'components/EntriesSelector/EntriesSelector';
import { ApplicationDataRow } from 'pages/Report/ApplicationSection/ApplicationSection';
import LinkComponent from 'components/LinkComponent/LinkComponent';
import KindTag from 'components/KindTag/KindTag.tsx';

interface ApplicationsEntriesSelectorProps {
    selectedApplications: ApplicationDataRow[];
    setSelectedApplications: React.Dispatch<React.SetStateAction<ApplicationDataRow[]>>;
    applicationsToExclude: ApplicationDataRow[];
    onAdd: () => void;
    onClose: () => void;
    availableApplications: ApplicationDataRow[];
}

const ApplicationsEntriesSelector: React.FC<ApplicationsEntriesSelectorProps> = ({
                                                                    selectedApplications,
                                                                    setSelectedApplications,
                                                                    applicationsToExclude,
                                                                    onAdd,
                                                                    onClose,
                                                                    availableApplications,
                                                                                 }) => {

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
        {
            header: 'Kind',
            columnKey: 'kind',
            customComponent: (app: ApplicationDataRow) => (
                <KindTag
                    name={app.kind || 'unknown'}
                />
            ),
        },
    ];

    return (
        <EntriesSelector<ApplicationDataRow>
            selectedItems={selectedApplications}
            setSelectedItems={setSelectedApplications}
            itemsToExclude={applicationsToExclude}
            onAdd={onAdd}
            onClose={onClose}
            items={availableApplications}
            columns={columns}
            getKey={getUniqueKey}
            entityLabel="application"
            noEntriesMessage={<p>There is no application to add.</p>}
            title="Select Applications"
        />
    );
};

export default ApplicationsEntriesSelector;

import './ApplicationsEntriesSelector.scss';
import React, {useState, useEffect} from 'react';
import Table, {TableColumn} from 'components/Table/Table';
import Checkbox from 'components/Checkbox/Checkbox';
import {ManagmentServiceApiInstance, AccuracyLevel} from 'api/managment-service';
import LinkComponent from 'components/LinkComponent/LinkComponent';
import {ApplicationDataRow} from 'pages/Report/ApplicationSection/ApplicationSection';
import ActionButton, {ActionButtonColor} from 'components/ActionButton/ActionButton.tsx';

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
    const [applications, setApplications] = useState<ApplicationDataRow[]>([]);
    const [selectAll, setSelectAll] = useState<boolean>(false);

    useEffect(() => {
        const fetchApplications = async () => {
            try {
                const data = await ManagmentServiceApiInstance.getApplications(clusterId);
                const rows = data.map((app) => ({
                    name: app.name,
                    running: app.running,
                    kind: app.kind,
                    accuracy: defaultAccuracy,
                    customPrompt: '',
                }));
                setApplications(rows);
            } catch (error) {
                console.error('Failed to fetch applications:', error);
            }
        };
        fetchApplications();
    }, [clusterId]);

    const availableApplications = applications.filter(
        (app) =>
            !applicationsToExclude.some(
                (excluded) => excluded.name === app.name
            )
    );

    useEffect(() => {
        setSelectAll(
            availableApplications.length > 0 &&
            selectedApplications.length === availableApplications.length
        );
    }, [selectedApplications, availableApplications]);

    const handleSelectAllChange = () => {
        if (selectAll) {
            setSelectedApplications([]);
        } else {
            setSelectedApplications(availableApplications);
        }
        setSelectAll(!selectAll);
    };

    const handleCheckboxChange = (app: ApplicationDataRow) => {
        setSelectedApplications((prevSelected) => {
            const isSelected = prevSelected.some(
                (selectedApp) => selectedApp.name === app.name
            );
            return isSelected
                // eslint-disable-next-line max-len
                ? prevSelected.filter((selectedApp) => selectedApp.name !== app.name)
                : [...prevSelected, app];
        });
    };

    const columns: TableColumn<ApplicationDataRow>[] = [
        {
            header: (
                <Checkbox
                    checked={selectAll}
                    onChange={handleSelectAllChange}
                />
            ),
            columnKey: 'checkbox',
            customComponent: (row: ApplicationDataRow) => (
                <Checkbox
                    checked={selectedApplications.some(
                        (selectedApp) => selectedApp.name === row.name
                    )}
                    onChange={() => handleCheckboxChange(row)}
                />
            ),
        },
        {
            header: 'Name',
            columnKey: 'name',
            customComponent: (row: ApplicationDataRow) => (
                <LinkComponent to="#" isRunning={row.running}>
                    {row.name}
                </LinkComponent>
            ),
        },
        {header: 'Kind', columnKey: 'kind'},
    ];

    return (
        <div className="application-entries">
            <Table
                columns={columns}
                rows={availableApplications.map((app) => ({
                    ...app,
                    key: app.name,
                }))}
            />
            <div className="application-entries__buttons">
                <ActionButton
                    onClick={onAdd}
                    description="Add"
                    color={ActionButtonColor.GREEN}
                />
                <ActionButton
                    onClick={onClose}
                    description="Close"
                    color={ActionButtonColor.RED}
                />
            </div>
        </div>
    );
};

export default ApplicationsEntriesSelector;

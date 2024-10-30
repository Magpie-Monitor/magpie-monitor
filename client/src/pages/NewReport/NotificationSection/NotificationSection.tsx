import './NotificationSection.scss';
import SectionComponent from 'components/SectionComponent/SectionComponent.tsx';
import Table, { TableColumn } from 'components/Table/Table.tsx';
import { useState, useEffect } from 'react';
// import ActionButton, { ActionButtonColor } from 'components/ActionButton/ActionButton.tsx';
import { NotificationChannel } from 'api/managment-service';

const MOCK_CHANNELS: NotificationChannel[] = [
    { id: '1', name: 'Infra team slack', service: 'Slack', details: 'wms_dev/#infra-alerts', added: '07.03.2024 15:32' },
    { id: '2', name: 'Infra team discord', service: 'Discord', details: 'wms_dev/#dev-infra-alerts', added: '07.03.2024 15:32' },
];

const NotificationSection = () => {
    const [rows, setRows] = useState<NotificationChannel[]>([]);
    const [loading, setLoading] = useState<boolean>(true);

    // const handleEdit = (channel: NotificationChannel) => {
    //     console.log('Edit channel:', channel);
    // };
    //
    // const handleDelete = (id: string) => {
    //     console.log('Delete channel with id:', id);
    //     setRows((prevRows) => prevRows.filter((row) => row.id !== id));
    // };

    const columns: Array<TableColumn<NotificationChannel>> = [
        { header: 'Name', columnKey: 'name' },
        { header: 'Service', columnKey: 'service' },
        { header: 'Details', columnKey: 'details' },
        { header: 'Added', columnKey: 'added' },
        {
            header: 'Actions',
            columnKey: 'actions',
            // render: (channel: NotificationChannel) => (
            //     <>
            //         <ActionButton onClick={() => handleEdit(channel)} description="Edit" color={ActionButtonColor.GREEN} />
            //         <ActionButton onClick={() => handleDelete(channel.id)} description="Delete" color={ActionButtonColor.RED} />
            //     </>
            // ),
        },
    ];

    useEffect(() => {
        setLoading(true);
        setTimeout(() => {
            setRows(MOCK_CHANNELS);
            setLoading(false);
        }, 2000);
    }, []);

    return (
        <SectionComponent icon={'notification-icon'} title={'Notification Channels'}>
            <div className="notification-section__content">
                {loading ? (
                    <p>Loading...</p>
                ) : rows.length === 0 ? (
                    <p>No notification channels available.</p>
                ) : (
                    <Table columns={columns} rows={rows} />
                )}
            </div>
        </SectionComponent>
    );
};

export default NotificationSection;

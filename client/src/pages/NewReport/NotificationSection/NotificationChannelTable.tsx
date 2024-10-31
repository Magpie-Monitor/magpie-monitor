import React from 'react';
import Table, { TableColumn } from 'components/Table/Table';
import Channels from './NotificationChannelsColumn/NotificationChannelColumn';
import ActionButton, { ActionButtonColor } from 'components/ActionButton/ActionButton';
import { NotificationChannel } from 'api/managment-service';
import {
    transformNotificationChannelToServiceColumn,
    transformNotificationChannelToDetailsColumn
} from './NotificationUtils.tsx';

interface NotificationChannelTableProps {
    rows: NotificationChannel[];
    onDelete: (id: string) => void;
}

const NotificationChannelTable: React.FC<NotificationChannelTableProps> = ({ rows, onDelete }) => {
    const columns: Array<TableColumn<NotificationChannel>> = [
        { header: 'Name', columnKey: 'name' },
        {
            header: 'Service',
            columnKey: 'service',
            customComponent: (row) => (
                <Channels channel={transformNotificationChannelToServiceColumn(row)} />
            ),
        },
        {
            header: 'Details',
            columnKey: 'details',
            customComponent: (row) => (
                <Channels channel={transformNotificationChannelToDetailsColumn(row)} />
            ),
        },
        { header: 'Updated', columnKey: 'updated' },
        { header: 'Added', columnKey: 'added' },
        {
            header: 'Actions',
            columnKey: 'actions',
            customComponent: (row: NotificationChannel) => (
                <ActionButton
                    onClick={() => onDelete(row.id)}
                    description="Delete"
                    color={ActionButtonColor.RED}
                />
            ),
        },
    ];

    return <Table columns={columns} rows={rows} />;
};

export default NotificationChannelTable;

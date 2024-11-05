import React from 'react';
import Table, { TableColumn } from 'components/Table/Table';
import Channels from './NotificationChannelsColumn/NotificationChannelColumn';
import ActionButton, {
    ActionButtonColor,
} from 'components/ActionButton/ActionButton';
import { NotificationChannel } from './NotificationSection.tsx';
import {
    transformNotificationChannelToServiceColumn,
    transformNotificationChannelToDetailsColumn,
} from './NotificationUtils.tsx';
import LinkComponent from 'components/LinkComponent/LinkComponent.tsx';

interface NotificationChannelTableProps {
    rows: NotificationChannel[];
    onDelete: (id: string, service: string) => void;
}

const NotificationChannelTable: React.FC<NotificationChannelTableProps> = ({
    rows,
    onDelete,
}) => {
    const columns: Array<TableColumn<NotificationChannel>> = [
        {
            header: 'Name',
            columnKey: 'name',
            customComponent: (row) => (
                <LinkComponent to="#">{row.name}</LinkComponent>
            ),
        },
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
        { header: 'Updated at', columnKey: 'updated' },
        { header: 'Added at', columnKey: 'added' },
        {
            header: 'Actions',
            columnKey: 'actions',
            customComponent: (row: NotificationChannel) => (
                <ActionButton
                    onClick={() => onDelete(row.id, row.service)}
                    description="Delete"
                    color={ActionButtonColor.RED}
                />
            ),
        },
    ];

    return <Table columns={columns} rows={rows} />;
};

export default NotificationChannelTable;

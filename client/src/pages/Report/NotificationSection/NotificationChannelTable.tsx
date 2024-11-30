import React from 'react';
import Table, { TableColumn } from 'components/Table/Table';
import Channels from './NotificationChannelsColumn/NotificationChannelColumn';
import { NotificationChannel } from './NotificationSection.tsx';
import {
    transformNotificationChannelToServiceColumn,
    transformNotificationChannelToDetailsColumn,
} from './NotificationUtils.tsx';
import LinkComponent from 'components/LinkComponent/LinkComponent.tsx';
import DeleteIconButton from 'components/DeleteIconButton/DeleteIconButton.tsx';

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
                <LinkComponent>{row.name}</LinkComponent>
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
            customComponent: (row) => (
                <DeleteIconButton onClick={() => onDelete(row.id, row.service)} />
            ),
        },
    ];

    return <Table columns={columns} rows={rows} />;
};

export default NotificationChannelTable;

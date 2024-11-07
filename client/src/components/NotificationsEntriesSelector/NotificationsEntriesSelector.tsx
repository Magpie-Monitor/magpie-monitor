import './NotificationEntriesSelector.scss';
import React, {useState, useEffect} from 'react';
import Table, {TableColumn} from 'components/Table/Table';
import {ManagmentServiceApiInstance} from 'api/managment-service.ts';
import {NotificationChannel} from 'pages/Report/NotificationSection/NotificationSection';
import Checkbox from 'components/Checkbox/Checkbox';
import Channels from 'pages/Report/NotificationSection/NotificationChannelsColumn/NotificationChannelColumn.tsx';
import {
    transformNotificationChannelToDetailsColumn,
    transformNotificationChannelToServiceColumn
} from 'pages/Report/NotificationSection/NotificationUtils.tsx';
import LinkComponent from 'components/LinkComponent/LinkComponent.tsx';
import ActionButton, {ActionButtonColor} from 'components/ActionButton/ActionButton.tsx';

interface EntriesSelectorProps {
    selectedChannels: NotificationChannel[];
    setSelectedChannels: React.Dispatch<React.SetStateAction<NotificationChannel[]>>;
    channelsToExclude: NotificationChannel[];
    onAdd: () => void;
    onClose: () => void;
}

const NotificationsEntriesSelector: React.FC<EntriesSelectorProps> = ({
                                                                          selectedChannels,
                                                                          setSelectedChannels,
                                                                          channelsToExclude,
                                                                          onAdd,
                                                                          onClose,
                                                                      }) => {
    const [notificationChannels, setNotificationChannels] = useState<NotificationChannel[]>([]);
    const [selectAll, setSelectAll] = useState<boolean>(false);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const channels = await ManagmentServiceApiInstance.getNotificationChannels();
                const channelRows = channels.map((channel): NotificationChannel => ({
                    id: channel.id,
                    name: channel.name,
                    service: channel.service,
                    details: channel.details,
                    updated: channel.updated,
                    added: channel.added,
                }));
                setNotificationChannels(channelRows);
            } catch (error) {
                console.error('Failed to fetch channels:', error);
            }
        };

        fetchData();
    }, []);

    const availableChannels = notificationChannels.filter(
        (channel) => !channelsToExclude.some(
            (existing) => existing.id === channel.id && existing.service === channel.service
        )
    );

    useEffect(() => {
        setSelectAll(availableChannels.length > 0
            && selectedChannels.length === availableChannels.length);
    }, [selectedChannels, availableChannels]);

    const getUniqueKey = (channel: NotificationChannel) => `${channel.id}-${channel.service}`;

    const handleSelectAllChange = () => {
        if (selectAll) {
            setSelectedChannels([]);
        } else {
            setSelectedChannels(availableChannels);
        }
        setSelectAll(!selectAll);
    };

    const handleCheckboxChange = (channel: NotificationChannel) => {
        setSelectedChannels((prevSelected: NotificationChannel[]) => {
            const isSelected = prevSelected.some(
                (selectedChannel) => getUniqueKey(selectedChannel) === getUniqueKey(channel)
            );
            return isSelected
                ? prevSelected.filter(
                    (selectedChannel) => getUniqueKey(selectedChannel) !== getUniqueKey(channel)
                )
                : [...prevSelected, channel];
        });
    };

    const columns: TableColumn<NotificationChannel>[] = [
        {
            header: (
                <Checkbox
                    checked={selectAll}
                    onChange={handleSelectAllChange}
                />
            ),
            columnKey: 'id-service',
            customComponent: (row: NotificationChannel) => (
                <Checkbox
                    checked={selectedChannels.some(
                        (selectedChannel) => getUniqueKey(selectedChannel) === getUniqueKey(row)
                    )}
                    onChange={() => handleCheckboxChange(row)}
                />
            ),
        },
        {
            header: 'Name',
            columnKey: 'name',
            customComponent: (row) => (
                <LinkComponent to="#">
                    {row.name}
                </LinkComponent>
            ),
        },
        {
            header: 'Service',
            columnKey: 'service',
            customComponent: (row) => (
                <Channels channel={transformNotificationChannelToServiceColumn(row)}/>
            ),
        },
        {
            header: 'Details',
            columnKey: 'details',
            customComponent: (row) => (
                <Channels channel={transformNotificationChannelToDetailsColumn(row)}/>
            ),
        },
    ];

    return (
        <div className="notification-entries">
            {availableChannels.length === 0 ? (
                <div className="notification-entries__no-channels-message">
                    <p>
                        No notification channels to display.
                    </p>
                    <LinkComponent to="/settings"
                                   className="notification-entries__no-channels-message__link">
                        You can create new one here.
                    </LinkComponent>
                </div>
            ) : (
                <Table
                    columns={columns}
                    rows={availableChannels.map(channel => ({
                        ...channel,
                        key: getUniqueKey(channel)
                    }))}
                    alignLeft={false}
                />
            )}
            <div className="notification-entries__buttons">
                {availableChannels.length > 0 && (
                    <ActionButton
                        onClick={onAdd}
                        description="Add"
                        color={ActionButtonColor.GREEN}
                    />
                )}
                <ActionButton
                    onClick={onClose}
                    description="Close"
                    color={ActionButtonColor.RED}
                />
            </div>
        </div>
    );
};

export default NotificationsEntriesSelector;

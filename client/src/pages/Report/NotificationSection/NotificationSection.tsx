import SectionComponent from 'components/SectionComponent/SectionComponent';
import SVGIcon from 'components/SVGIcon/SVGIcon';
import {useEffect, useState} from 'react';
import NotificationChannelTable from './NotificationChannelTable';
import OverlayComponent from 'components/OverlayComponent/OverlayComponent.tsx';
import NotificationsEntriesSelector
    from 'components/NotificationsEntriesSelector/NotificationsEntriesSelector.tsx';

export interface NotificationChannel {
    id: string;
    name: string;
    service: string;
    details: string;
    updated: string;
    added: string;

    [key: string]: string;
}

interface NotificationSectionProps {
    setNotificationChannels: (channels: NotificationChannel[]) => void;
}

const NotificationSection: React.FC<NotificationSectionProps> = ({setNotificationChannels}) => {
    const [rows, setRows] = useState<NotificationChannel[]>([]);
    const [selectedChannels, setSelectedChannels] = useState<NotificationChannel[]>([]);
    const [showModal, setShowModal] = useState(false);

    useEffect(() => {
        setNotificationChannels(rows);
    }, [rows, setNotificationChannels]);

    const handleAddClick = () => {
        setShowModal(true);
    };

    const handleCloseModal = () => {
        setShowModal(false);
    };

    const handleDelete = (id: string, service: string) => {
        setRows((prevRows) => {
            return prevRows.filter((row) => !(row.id === id && row.service === service));
        });
    };

    const handleAddSelected = () => {
        setRows((prevRows) => {
            const newChannels = selectedChannels.filter(
                (channel) => !prevRows.some(
                    (row) => row.id === channel.id && row.service === channel.service
                )
            );
            return [...prevRows, ...newChannels];
        });
        setShowModal(false);
        setSelectedChannels([]);
    };

    return (
        <SectionComponent
            icon={<SVGIcon iconName="notification-icon"/>}
            title="Notification channels"
            callback={handleAddClick}
        >
            <OverlayComponent
                isDisplayed={showModal}
                onClose={handleCloseModal}
            >
                <NotificationsEntriesSelector
                    selectedChannels={selectedChannels}
                    setSelectedChannels={setSelectedChannels}
                    channelsToExclude={rows}
                    onAdd={handleAddSelected}
                    onClose={handleCloseModal}
                />
            </OverlayComponent>
            {rows.length === 0 ? (
                <p>No notification channels selected, please add new.</p>
            ) : (
                <NotificationChannelTable rows={rows} onDelete={handleDelete}/>
            )}
        </SectionComponent>
    );
};

export default NotificationSection;
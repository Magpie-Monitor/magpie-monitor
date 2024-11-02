import SectionComponent from 'components/SectionComponent/SectionComponent';
import SVGIcon from 'components/SVGIcon/SVGIcon';
import { useState, useEffect } from 'react';
import NotificationChannelTable from './NotificationChannelTable';
import OverlayComponent from 'components/OverlayComponent/OverlayComponent.tsx';
import { ManagmentServiceApiInstance } from 'api/managment-service';
import Spinner from 'components/Spinner/Spinner.tsx';

export interface NotificationChannel {
    id: string;
    name: string;
    service: string;
    details: string;
    updated: string;
    added: string;
    [key: string]: string;
}

const NotificationSection = () => {
    const [rows, setRows] = useState<NotificationChannel[]>([]);
    const [loading, setLoading] = useState(true);
    const [showModal, setShowModal] = useState(false);

    const fetchNotificationChannels = async () => {
        try {
            setLoading(true);
            const channelsData = await ManagmentServiceApiInstance.getNotificationChannels();  // Assuming getNotificationChannels() fetches notification channels

            const channelRows = channelsData.map((channel): NotificationChannel => ({
                id: channel.id,
                name: channel.name,
                service: channel.service,
                details: channel.details,
                updated: channel.updated,
                added: channel.added,
            }));

            setRows(channelRows);
        } catch (e: unknown) {
            console.error('Failed to fetch notification channels', e);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchNotificationChannels();
    }, []);

    const handleAddClick = () => {
        setShowModal(true);
    };

    const handleCloseModal = () => {
        setShowModal(false);
    };

    const handleDelete = (id: string) => {
        setRows((prevRows) => prevRows.filter((row) => row.id !== id));
    };

    return (
        <SectionComponent
            icon={<SVGIcon iconName="notification-icon" />}
            title="Notification channels"
            callback={handleAddClick}
        >
            {showModal && (
                <OverlayComponent isDisplayed={showModal} onClose={handleCloseModal}>
                    <p>No notification channels here (probably Wojciech dropped all of them)</p>
                </OverlayComponent>
            )}
            {loading ? (
                <Spinner />
            ) : rows.length === 0 ? (
                <p>No notification channels selected, please add new.</p>
            ) : (
                <NotificationChannelTable rows={rows} onDelete={handleDelete} />
            )}
        </SectionComponent>
    );
};

export default NotificationSection;

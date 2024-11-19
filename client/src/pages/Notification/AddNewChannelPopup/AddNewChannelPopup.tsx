export interface AddNewChannelPopupProps {
  isDisplayed: boolean;
  setIsDisplayed: (arg: boolean) => void;
}

export interface AddNewChannelPopup {
  popup: ({
    isDisplayed,
    setIsDisplayed,
  }: AddNewChannelPopupProps) => React.ReactNode;
}

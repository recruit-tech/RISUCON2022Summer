package rcalendar.model;

import org.jboss.resteasy.reactive.multipart.FileUpload;

import javax.ws.rs.FormParam;

public class UpdateIconRequest {
        @FormParam("icon")
        private FileUpload icon;

        public FileUpload getIcon() {
                return icon;
        }

        public void setIcon(FileUpload icon) {
                this.icon = icon;
        }
}

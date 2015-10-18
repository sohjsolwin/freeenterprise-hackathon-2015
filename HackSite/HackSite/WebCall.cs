using System;
using System.Net;
using System.Text;

namespace HackSite
{
    public class WebCall
    {
        /// <summary>
        /// Gets or sets the API URI.
        /// </summary>
        /// <value>
        /// The API URI.
        /// </value>
        public string APIUri { get; set; }

        /// <summary>
        /// Gets or sets the headers.
        /// </summary>
        /// <value>
        /// The headers.
        /// </value>
        protected WebHeaderCollection Headers { get; set; }

        /// <summary>
        /// Initializes a new instance of the <see cref="WebAPIInterface"/> class.
        /// </summary>
        public WebCall()
        {
            this.Headers = new WebHeaderCollection();
        }

        /// <summary>
        /// Posts the specified post style.
        /// </summary>
        /// <param name="postData">The post style.</param>
        /// <returns></returns>
        internal string PostString(string postData)
        {
            string returnData = null;

            using (WebClient webclient = new WebClient())
            {
                webclient.Headers = this.Headers;
                returnData = webclient.UploadString(this.APIUri, postData);
            }

            return returnData;
        }

        /// <summary>
        /// Posts the specified post style.
        /// </summary>
        /// <param name="postData">The post style.</param>
        /// <returns></returns>
        internal string PostData(string postData)
        {
            string returnString = null;

            using (WebClient webclient = new WebClient())
            {
                webclient.Headers = this.Headers;

                var data = Encoding.UTF8.GetBytes(postData);

                byte[] incomingData = webclient.UploadData(this.APIUri, data);
                returnString = Encoding.UTF8.GetString(incomingData);
            }
            return returnString;
        }

        /// <summary>
        /// Posts the specified post style.
        /// </summary>
        /// <param name="postData">The post style.</param>
        /// <returns></returns>
        internal byte[] PostMassData(string postData)
        {
            byte[] returnData;
            using (WebClient webclient = new WebClient())
            {
                webclient.Headers = this.Headers;

                var data = Encoding.UTF8.GetBytes(postData);

                returnData = webclient.UploadData(this.APIUri, data);
            }
            return returnData;
        }

        /// <summary>
        /// Sends the request. The resonse will come in Request_DownloadCompleted Event.
        /// </summary>
        internal string Get()
        {
            string returnString = null;
            using (WebClient webclient = new WebClient())
            {
                webclient.Headers = this.Headers;
                returnString = webclient.DownloadString(new Uri(this.APIUri));
            }
            return returnString;
        }
    }
}

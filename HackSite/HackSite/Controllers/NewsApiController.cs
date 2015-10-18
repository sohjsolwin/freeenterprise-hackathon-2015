using HackSite.Dto;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Web;
using System.Web.Http;

namespace HackSite.Controllers
{
    [RoutePrefix("api/newsdata")]
    public class NewsApiController : ApiController
    {

        private static IEnumerable<NewsData> cachedNewsData;


        [Route("getData")]
        public IHttpActionResult GetGraphData()
        {
            var startDate = new DateTime(DateTime.Today.Year, DateTime.Today.Month, 1);
            if (cachedNewsData == null || cachedNewsData.Count() <= 0)
            {
                cachedNewsData = ParseData.CallHavenAPI(startDate);
            }

            var data = cachedNewsData.Select(x => new { x.State, x.score }).Distinct()
                .ToDictionary(x => x.State, x => (object)((int)x.score));
            //var modData = new List<KeyValuePair<string, object>>() { (new KeyValuePair<string, object>("State", "Unemployment")) }.Union(data).ToDictionary();
            ;


            return Json(ConvertObjToSerializedArr(data));
        }

        [Route("getStateInfo")]
        public IHttpActionResult GetStateInfo([FromUri] string state)
        {
            return Json(cachedNewsData.Where(x => x.State == state).ToList());
        }

        private List<object[]> ConvertObjToSerializedArr(Dictionary<string, object> data)
        {
            List<object[]> dataSet = new List<object[]>();
            foreach (var dataPair in data)
            {
                var itemElement = new object[2] { dataPair.Key, dataPair.Value };
                dataSet.Add(itemElement);
            }

            dataSet.Insert(0, new object[] { "State", "Recent News Sentiment" });

            return dataSet;
        }

        

    }
}
